package bll

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"github.com/cloudslit/cloudslit/client/internal/config"
	"github.com/cloudslit/cloudslit/client/internal/contextx"
	"github.com/cloudslit/cloudslit/client/internal/schema"
	"github.com/cloudslit/cloudslit/client/pkg/certificate"
	"github.com/cloudslit/cloudslit/client/pkg/errors"
	"github.com/cloudslit/cloudslit/client/pkg/logger"
	"github.com/cloudslit/cloudslit/client/pkg/recover"
	"github.com/cloudslit/cloudslit/client/pkg/util/trace"
	"github.com/xtaci/smux"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (a *Client) DialWS(ctx context.Context, nextAddr *schema.NextServer) (net.Conn, error) {
	conn, err := tls.Dial("tcp", nextAddr.Host+":"+nextAddr.Port, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	secretLink := "secretLink"
	req, err := http.NewRequest("GET", "/"+secretLink, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Host = nextAddr.Host
	traceID, _ := contextx.FromTraceID(ctx)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("X-TraceID", traceID)
	req.Header.Set("X-ClientCert", base64.StdEncoding.EncodeToString([]byte(config.C.Certificate.CertPem)))
	req.Header.Set("X-ClientCA", base64.StdEncoding.EncodeToString([]byte(config.C.Certificate.CaPem)))

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
		// Obtain the server certificate
		serverCa := resp.Header.Get("X-ServerCert")
		if resp.Header.Get("X-ServerCert") == "" {
			err = errors.New("Failed to obtain the lower-layer service certificate. Procedure")
			return nil, errors.WithStack(err)
		}
		serverCaCert, err := base64.StdEncoding.DecodeString(serverCa)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		// Verify the server certificate
		err = certificate.NewVerify(string(serverCaCert), config.C.Certificate.CaPem, nextAddr.Host).Verify()
		if err != nil {
			//return nil, errors.WithStack(err)
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

func (a *Client) Listen(ctx context.Context) error {
	lisAddr := config.C.Common.LocalAddr + ":" + strconv.Itoa(config.C.Common.LocalPort)
	ln, err := net.Listen("tcp", lisAddr)
	if err != nil {
		return err
	}
	logger.WithContext(ctx).Printf("Started Client at %v\n", ln.Addr().String())

	for {
		clientConn, err := ln.Accept()
		if err != nil {
			logger.WithErrorStack(ctx, errors.WithStack(err)).Error("Failed to accept connection:", err)
			continue
		}
		recover.Recovery(ctx, func() {
			a.handleConn(ctx, clientConn)
		})
	}
}

func (a *Client) handleConn(ctx context.Context, clientConn net.Conn) {
	traceID := trace.NewTraceID()
	ctx = contextx.NewTraceID(ctx, traceID)
	ctx = logger.NewTraceIDContext(ctx, traceID)
	defer func() {
		closeErr := clientConn.Close()
		if closeErr != nil {
			logger.WithErrorStack(ctx, errors.WithStack(closeErr)).Errorf("Closed Connection with error: %v\n", closeErr)
		} else {
			logger.WithContext(ctx).Infof("Closed Connection: %v\n", clientConn.RemoteAddr().String())
		}
	}()
	// TODO: 获取请求服务地址
	//nextServer := a.GetNextServer(conf)
	nextServer := &schema.NextServer{
		Host: "https://",
		Port: "443",
	}
	serverConn, err := a.DialWS(ctx, nextServer)
	if err != nil {
		logger.WithErrorStack(ctx, err).Errorf("The client failed to request the lower level service. Procedure:Addr:%s:%s Error:%v", nextServer.Host, nextServer.Port, err)
		return
	}
	defer serverConn.Close()
	// 多路复用
	//Setup client side of smux
	session, err := smux.Client(serverConn, nil)
	if err != nil {
		logger.WithErrorStack(ctx, err).Errorf("smux.Client. Procedure:Addr:%s:%s Error:%v", nextServer.Host, nextServer.Port, err)
		return
	}
	defer session.Close()
	// Open a new stream
	stream, err := session.OpenStream()
	if err != nil {
		logger.WithErrorStack(ctx, err).Errorf("session.OpenStream. Procedure:Addr:%s:%s Error:%v", nextServer.Host, nextServer.Port, err)
		return
	}
	defer stream.Close()
	logger.WithContext(ctx).Infof("Connect Success:Addr:%s:%s", nextServer.Host, nextServer.Port)
	TransparentProxy(clientConn, stream)
}

func (a *Client) GetNextServer(chains *schema.ClientConfig) *schema.NextServer {
	replyCount := len(chains.Relays)
	nextServer := new(schema.NextServer)
	if replyCount == 0 {
		nextServer.Host = chains.Server.Host
		nextServer.Port = strconv.Itoa(chains.Server.OutPort)
	} else {
		chain := chains.Relays[0]
		nextServer.Host = chain.Host
		nextServer.Port = strconv.Itoa(chain.OutPort)
	}
	return nextServer
}
