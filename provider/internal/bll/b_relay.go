package bll

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/cloudslit/cloudslit/provider/internal/config"
	"github.com/cloudslit/cloudslit/provider/internal/contextx"
	"github.com/cloudslit/cloudslit/provider/internal/metrics"
	"github.com/cloudslit/cloudslit/provider/internal/schema"
	"github.com/cloudslit/cloudslit/provider/pkg/certificate"
	"github.com/cloudslit/cloudslit/provider/pkg/errors"
	"github.com/cloudslit/cloudslit/provider/pkg/logger"
	"github.com/cloudslit/cloudslit/provider/pkg/pconst"
	"github.com/cloudslit/cloudslit/provider/pkg/recover"
	"github.com/cloudslit/cloudslit/provider/pkg/util/json"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

type Relay struct{}

// ReadInitiaWSRequest Receiving WS Requests
func (a *Relay) ReadInitiaWSRequest(ctx context.Context, conf *schema.RelayConfig, connReader *bufio.Reader) (*schema.ClientConfig, *http.Request, context.Context, error) {
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
		chainsJSON := strings.ToLower(req.Header.Get("X-Chains"))
		if chainsJSON == "" {
			return nil, nil, ctx, errors.NewWithStack("X-Chains argument is missing")
		}
		var chains schema.ClientConfig
		err = json.Unmarshal([]byte(chainsJSON), &chains)
		if err != nil {
			return nil, nil, ctx, errors.WithStack(err)
		}
		// get client cert
		clientCa := req.Header.Get("X-ClientCert")
		if clientCa == "" {
			return nil, nil, ctx, errors.NewWithStack("X-ClientCert argument is missing")
		}
		clientCaCert, err := base64.StdEncoding.DecodeString(clientCa)
		if err != nil {
			return nil, nil, ctx, errors.WithStack(err)
		}
		// check client cert
		err = certificate.NewVerify(string(clientCaCert), config.C.Certificate.CaPem, "").Verify()
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
	err = errors.New("Reply Illegal request:\n" + string(reqBytes))
	return nil, nil, ctx, errors.WithStack(err)
}

// Responding to WS requests
func (a *Relay) GenerateInitialWSResponse(ctx context.Context, clientConn net.Conn, req *http.Request) ([]byte, error) {
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
	resp.Header.Set("X-ServerCert", base64.StdEncoding.EncodeToString([]byte(config.C.Certificate.CertPem)))

	res, err := httputil.DumpResponse(&resp, true)
	_, err = clientConn.Write(res)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, err
}

func (a *Relay) handleConn(ctx context.Context, conf *schema.RelayConfig, clientConn net.Conn) error {
	begin := time.Now()
	defer func() {
		closeErr := clientConn.Close()
		if closeErr != nil {
			logger.WithErrorStack(ctx, errors.WithStack(closeErr)).Errorf("Closed Connection with error: %v\n", closeErr)
		} else {
			logger.WithContext(ctx).Infof("Closed Connection: %v\n", clientConn.RemoteAddr().String())
		}
	}()
	connReader := bufio.NewReader(clientConn)
	chains, req, ctx, err := a.ReadInitiaWSRequest(ctx, conf, connReader)
	if err != nil {
		metrics.AddDelayPoint(ctx, pconst.OperatorRelay, metrics.ReqFail, time.Now().Sub(begin).String(), conf.UUID, conf.Name)
		logger.WithErrorStack(ctx, err).Error("Error obtaining WS request information：", err)
		return err
	}
	_, err = a.GenerateInitialWSResponse(ctx, clientConn, req)
	if err != nil {
		metrics.AddDelayPoint(ctx, pconst.OperatorRelay, metrics.ReqFail, time.Now().Sub(begin).String(), conf.UUID, conf.Name)
		logger.WithErrorStack(ctx, err).Error("Response WS message error：", err)
		return err
	}
	// Get server certificate verification information
	verifyFlag := "serverCaReady"
	verifyBytes, err := connReader.Peek(len(verifyFlag))
	if err != nil {
		metrics.AddDelayPoint(ctx, pconst.OperatorRelay, metrics.ReqFail, time.Now().Sub(begin).String(), conf.UUID, conf.Name)
		logger.WithErrorStack(ctx, err).Error("Error obtaining the certificate verification result：", err)
		return err
	}
	// Server certificate verification passed
	if string(verifyBytes) == verifyFlag {
		nextServer := a.GetNextServer(conf, chains)
		serverConn, err := a.DialWS(ctx, nextServer, req, conf, chains)
		if err != nil {
			metrics.AddDelayPoint(ctx, pconst.OperatorRelay, metrics.ReqFail, time.Now().Sub(begin).String(), conf.UUID, conf.Name)
			logger.WithErrorStack(ctx, err).Errorf("The relay side failed to request the lower-level service:Addr:%s:%s Error:%v", nextServer.Host, nextServer.Port, err)
			return err
		}
		defer serverConn.Close()
		end := time.Now().Sub(begin).String()
		logger.WithContext(ctx).Infof("Connect Success:Addr:%s:%s", nextServer.Host, nextServer.Port)
		metrics.AddDelayPoint(ctx, pconst.OperatorRelay, metrics.ReqSuccess, end, conf.UUID, conf.Name)
		TransparentProxy(clientConn, serverConn)
		return nil
	}
	err = errors.New("Relay side certificate verification failed\n")
	logger.WithErrorStack(ctx, errors.WithStack(err)).Error(err)
	metrics.AddDelayPoint(ctx, pconst.OperatorRelay, metrics.ReqFail, time.Now().Sub(begin).String(), conf.UUID, conf.Name)
	return err
}

func (a *Relay) DialWS(ctx context.Context, nextChain *schema.NextServer, req *http.Request, conf *schema.RelayConfig, chains *schema.ClientConfig) (net.Conn, error) {
	conn, err := tls.Dial("tcp", nextChain.Host+":"+nextChain.Port, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Host = nextChain.Host
	req.Header.Set("X-ClientCert", base64.StdEncoding.EncodeToString([]byte(config.C.Certificate.CertPem)))

	err = req.Write(conn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.Status == "101 Switching Protocols" &&
		strings.ToLower(resp.Header.Get("Upgrade")) == "websocket" &&
		strings.ToLower(resp.Header.Get("Connection")) == "upgrade" {
		// Get server certificate
		serverCa := resp.Header.Get("X-ServerCert")
		if resp.Header.Get("X-ServerCert") == "" {
			err = errors.New("Failed to obtain lower-level service certificate")
			return nil, errors.WithStack(err)
		}
		serverCaCert, err := base64.StdEncoding.DecodeString(serverCa)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		// Verify server certificate
		err = certificate.NewVerify(string(serverCaCert), config.C.Certificate.CaPem, nextChain.Host).Verify()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		_, err = conn.Write([]byte("serverCaReady"))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return conn, nil
	}
	respBytes, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = errors.New("Got unexpected response:\n" + string(respBytes) + "\nstatus:" + resp.Status)
	return nil, errors.WithStack(err)
}

func NewRelay() *Relay {
	return &Relay{}
}

func (a *Relay) Listen(ctx context.Context, attrs map[string]interface{}) {
	go func() {
		conf, err := schema.ParseRelayConfig(attrs)
		if err != nil {
			panic(err)
		}
		cert, err := tls.X509KeyPair([]byte(config.C.Certificate.CertPem), []byte(config.C.Certificate.KeyPem))
		if err != nil {
			panic(err)
		}
		l, err := tls.Listen("tcp", "0.0.0.0:"+strconv.Itoa(conf.Port), &tls.Config{
			Certificates: []tls.Certificate{cert},
		})
		if err != nil {
			panic(err)
		}
		logger.WithContext(ctx).Printf("Started ZERO ACCESS Relay at %v\n", l.Addr().String())
		for {
			conn, err := l.Accept()
			if err != nil {
				logger.WithErrorStack(ctx, errors.WithStack(err)).Error("Failed to accept connection:", err)
				continue
			}
			recover.Recovery(ctx, func() {
				a.handleConn(ctx, conf, conn)
			})
		}
	}()
}

func (a *Relay) GetNextServer(conf *schema.RelayConfig, chains *schema.ClientConfig) *schema.NextServer {
	replyCount := len(chains.Relays)
	nextKey := 0
	for key, item := range chains.Relays {
		// last
		if key == replyCount-1 {
			break
		}
		if item.UUID == conf.UUID {
			nextKey = key + 1
		}
	}
	nextServer := new(schema.NextServer)
	if nextKey == 0 {
		nextServer.Host = chains.Server.Host
		nextServer.Port = strconv.Itoa(chains.Server.OutPort)
	} else {
		chain := chains.Relays[nextKey]
		nextServer.Host = chain.Host
		nextServer.Port = strconv.Itoa(chain.OutPort)
	}
	return nextServer
}
