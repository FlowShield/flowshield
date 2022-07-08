package bll

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/cloudslit/cloudslit/provider/internal/contextx"
	"github.com/cloudslit/cloudslit/provider/internal/schema"
	"github.com/cloudslit/cloudslit/provider/pkg/certificate"
	"github.com/cloudslit/cloudslit/provider/pkg/errors"
	"github.com/cloudslit/cloudslit/provider/pkg/logger"
	"github.com/cloudslit/cloudslit/provider/pkg/recover"
	"github.com/cloudslit/cloudslit/provider/pkg/util/json"
	"github.com/xtaci/smux"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

const RootCA = ``

type Provider struct {
}

// ReadInitiaWSRequest Read the WS request
func (a *Provider) ReadInitiaWSRequest(ctx context.Context, connReader *bufio.Reader) (*schema.ClientConfig, *http.Request, context.Context, error) {
	expectedH1Req := "GET /secretLink"
	firstBytes, err := connReader.Peek(len(expectedH1Req))
	if err != nil {
		return nil, nil, ctx, errors.WithStack(err)
	}
	if string(firstBytes) == expectedH1Req {
		req, err := http.ReadRequest(connReader)
		if err != nil {
			return nil, nil, ctx, errors.WithStack(err)
		}
		traceID := req.Header.Get("X-TraceID")
		if traceID != "" {
			ctx = contextx.NewTraceID(ctx, traceID)
			ctx = logger.NewTraceIDContext(ctx, traceID)
		}
		if strings.ToLower(req.Header.Get("Connection")) != "upgrade" && strings.ToLower(req.Header.Get("Connection")) != "keep-alive, upgrade" {
			return nil, nil, ctx, errors.WithStack(fmt.Errorf("Connection header expected: upgrade, got: %s\n",
				strings.ToLower(req.Header.Get("Connection"))))
		}
		if strings.ToLower(req.Header.Get("Upgrade")) != "websocket" {
			return nil, nil, ctx, errors.WithStack(fmt.Errorf("Upgrade header expected: websocket, got: %s\n",
				strings.ToLower(req.Header.Get("Upgrade"))))
		}
		// Get link information
		chainsJSON := strings.ToLower(req.Header.Get("X-Chains"))
		if chainsJSON == "" {
			return nil, nil, ctx, errors.NewWithStack("X-Chains argument is missing")
		}
		var chains schema.ClientConfig
		err = json.Unmarshal([]byte(chainsJSON), &chains)
		if err != nil {
			return nil, nil, ctx, errors.WithStack(err)
		}
		// get client certificate
		clientCa := req.Header.Get("X-ClientCert")
		if clientCa == "" {
			return nil, nil, ctx, errors.NewWithStack("X-ClientCert argument is missing")
		}
		clientCaCert, err := base64.StdEncoding.DecodeString(clientCa)
		if err != nil {
			return nil, nil, ctx, errors.WithStack(err)
		}
		// Verify the client certificate
		err = certificate.NewVerify(string(clientCaCert), RootCA, "").Verify()
		if err != nil {
			return nil, nil, ctx, errors.WithStack(err)
		}
		return &chains, req, ctx, nil
	}
	req, err := http.ReadRequest(connReader)
	if err != nil {
		return nil, nil, ctx, errors.WithStack(err)
	}
	reqBytes, err := httputil.DumpRequest(req, false)
	if err != nil {
		return nil, nil, ctx, errors.WithStack(err)
	}
	err = errors.New("Server Illegal request:\n" + string(reqBytes))
	return nil, nil, ctx, errors.WithStack(err)
}

// Responding to WS requests
func (a *Provider) GenerateInitialWSResponse(ctx context.Context, clientConn net.Conn, req *http.Request) ([]byte, error) {
	resp := http.Response{
		Status:           "101 Switching Protocols",
		StatusCode:       101,
		Proto:            "HTTP/1.1",
		ProtoMajor:       1,
		ProtoMinor:       1,
		Header:           http.Header{},
		Body:             nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}
	resp.Header.Set("Upgrade", req.Header.Get("Upgrade"))
	resp.Header.Set("Connection", req.Header.Get("Connection"))
	resp.Header.Set("X-ServerCert", base64.StdEncoding.EncodeToString([]byte("config.C.Certificate.CertPem")))

	res, err := httputil.DumpResponse(&resp, true)
	_, err = clientConn.Write(res)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, err
}

func (a *Provider) handleConn(ctx context.Context, clientConn net.Conn) error {
	connReader := bufio.NewReader(clientConn)
	chains, req, ctx, err := a.ReadInitiaWSRequest(ctx, connReader)
	if err != nil {
		logger.WithErrorStack(ctx, err).Error("Error obtaining WS request information：", err)
		return err
	}
	_, err = a.GenerateInitialWSResponse(ctx, clientConn, req)
	if err != nil {
		logger.WithErrorStack(ctx, err).Error("Response WS message error：", err)
		return err
	}
	verifyFlag := "serverCaReady"
	verifyBytes, err := connReader.Peek(len(verifyFlag))
	if err != nil {
		logger.WithErrorStack(ctx, err).Error("Error obtaining the certificate verification result.：", err)
		return err
	}
	if string(verifyBytes) == verifyFlag {
		targetAddr := chains.Target.Host + ":" + strconv.Itoa(chains.Target.Port)
		serverConn, err := net.Dial("tcp", targetAddr)
		if err != nil {
			logger.WithErrorStack(ctx, errors.WithStack(err)).Errorf("Failed to request resource from server\n:Addr:%s Error:%v", targetAddr, err)
			return err
		}
		// 多路复用
		session, err := smux.Server(clientConn, nil)
		if err != nil {
			return err
		}
		stream, err := session.AcceptStream()
		if err != nil {
			return err
		}
		defer func() {
			closeErr := clientConn.Close()
			if closeErr != nil {
				logger.WithContext(ctx).Errorf("Closed Connection with error: %v", closeErr)
			} else {
				logger.WithContext(ctx).Infof("Closed Connection: %v", clientConn.RemoteAddr().String())
			}
			serverConn.Close()
			session.Close()
			stream.Close()
		}()
		logger.WithContext(ctx).Infof("Connection Success: Client:%s; To:%s;", clientConn.RemoteAddr().String(), targetAddr)
		TransparentProxy(stream, serverConn)
		return nil
	}
	err = errors.New("Server certificate verification failed")
	logger.WithErrorStack(ctx, errors.WithStack(err)).Error(err)
	return err
}

func NewProvider() *Provider {
	return &Provider{}
}

func (a *Provider) Listen(ctx context.Context, port int, content *schema.ProviderContent) (net.Listener, error) {
	cert, err := tls.X509KeyPair([]byte(content.CertPem), []byte(content.KeyPem))
	if err != nil {
		logger.Errorf("证书解析失败:%s", err)
		return nil, err
	}
	l, err := tls.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port), &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	if err != nil {
		logger.Errorf("监听端口失败:%s", err)
		return nil, err
	}
	logger.WithContext(ctx).Printf("Started Provider Server at %v", l.Addr().String())
	return l, nil
}

func (a *Provider) Handle(ctx context.Context, l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			logger.WithContext(ctx).Error("Failed to accept connection:", err)
			continue
		}
		recover.Recovery(ctx, func() {
			a.handleConn(ctx, conn)
		})
	}
}

func verifyPort(port int) (int, error) {
	var ln net.Listener
	var err error
	// Automatically look for an open port when a custom port isn't
	// selected by a user.
	for {
		ln, err = net.Listen("tcp", ":"+strconv.Itoa(port))
		if err == nil {
			break
		}
		if port >= 65535 {
			return port, errors.New("failed to find open port")
		}
		port++
	}
	if ln != nil {
		ln.Close()
	}
	return port, nil
}
