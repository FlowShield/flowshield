package bll

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/flowshield/flowshield/client/internal/config"
	"github.com/flowshield/flowshield/client/internal/contextx"
	"github.com/flowshield/flowshield/client/internal/schema"
	"github.com/flowshield/flowshield/client/pkg/errors"
	"github.com/flowshield/flowshield/client/pkg/logger"
	"github.com/flowshield/flowshield/client/pkg/recover"
	"github.com/flowshield/flowshield/client/pkg/util/trace"
	"github.com/xtaci/smux"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (a *Client) DialWS(ctx context.Context, client *schema.ClientConfig) (net.Conn, error) {
	serverAddr := client.Server.Host + ":" + strconv.Itoa(client.Server.Port)
	tlsc, err := a.GetDialMtlsConfig(client)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	//tlsc.ServerName = "cloud-slit"
	//serverAddr = "127.0.0.1:5092"
	conn, err := tls.Dial("tcp", serverAddr, tlsc)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	secretLink := "secretLink"
	req, err := http.NewRequest("GET", "/"+secretLink, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Host = client.Server.Host
	traceID, _ := contextx.FromTraceID(ctx)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("X-TraceID", traceID)
	req.Header.Set("X-Chains", client.String())

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
		return conn, nil
	}
	respBytes, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = errors.New("Got unexpected response:\n" + string(respBytes) + "\nstatus:" + resp.Status)
	return nil, errors.WithStack(err)
}

func (a *Client) GetDialMtlsConfig(client *schema.ClientConfig) (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(client.CertPem), []byte(client.KeyPem))
	if err != nil {
		return nil, err
	}
	// Load ca certificate
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(client.CaPem))
	tlsc := &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{cert}, // client certificate
		ServerName:   client.Server.Host,      // Server certificate common name
		RootCAs:      pool,                    // The server certificate belongs to ca
	}
	return tlsc, err
}

func (a *Client) Listen(ctx context.Context, client *schema.ClientConfig) error {
	lisAddr := config.C.App.LocalAddr + ":" + strconv.Itoa(config.C.App.LocalPort)
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
			a.handleConn(ctx, clientConn, client)
		})
	}
}

func (a *Client) handleConn(ctx context.Context, clientConn net.Conn, client *schema.ClientConfig) {
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
	serverConn, err := a.DialWS(ctx, client)
	if err != nil {
		logger.WithErrorStack(ctx, err).Errorf("The client failed to request the lower level service. Procedure:Addr:%s:%s Error:%v", client.Server.Host, client.Server.Port, err)
		return
	}
	defer serverConn.Close()
	// Multiplexing
	//Setup client side of smux
	session, err := smux.Client(serverConn, nil)
	if err != nil {
		logger.WithErrorStack(ctx, err).Errorf("smux.Client. Procedure:Addr:%s:%s Error:%v", client.Server.Host, client.Server.Port, err)
		return
	}
	defer session.Close()
	// Open a new stream
	stream, err := session.OpenStream()
	if err != nil {
		logger.WithErrorStack(ctx, err).Errorf("session.OpenStream. Procedure:Addr:%s:%s Error:%v", client.Server.Host, client.Server.Port, err)
		return
	}
	defer stream.Close()
	logger.WithContext(ctx).Infof("Connect Success:Addr:%s:%d", client.Server.Host, client.Server.Port)
	TransparentProxy(clientConn, stream)
}
