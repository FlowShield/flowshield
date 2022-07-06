package bll

import (
	"io"
	"net"
)

func TransparentProxy(clientConn, serverConn net.Conn) {
	errChan := make(chan error, 2)
	copyConn := func(a, b net.Conn) {
		_, err := io.Copy(a, b)
		errChan <- err
	}
	go copyConn(clientConn, serverConn)
	go copyConn(serverConn, clientConn)
	select {
	case <-errChan:
		return
	}
}
