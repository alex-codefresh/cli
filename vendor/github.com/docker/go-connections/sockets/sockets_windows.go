package sockets

import (
	"net"
	"net/http"
	"time"

	winio "github.com/Microsoft/go-winio"
	"github.com/sirupsen/logrus"
)

func configureUnixTransport(tr *http.Transport, proto, addr string) error {
	addr = "C:" + addr
	logrus.Warn("Socket address is: %s", addr)
	addrUnix, _ := net.ResolveUnixAddr("unix", addr)
	tr.DisableCompression = true
	tr.Dial = func(_, _ string) (net.Conn, error) {
		c, err := net.DialUnix(proto, nil, addrUnix)
		if err != nil {
			logrus.Errorf("Unix connection failed: %v", err)
		}
		return c, err
	}
	return nil
}

func configureNpipeTransport(tr *http.Transport, proto, addr string) error {
	// No need for compression in local communications.
	tr.DisableCompression = true
	tr.Dial = func(_, _ string) (net.Conn, error) {
		return DialPipe(addr, defaultTimeout)
	}
	return nil
}

// DialPipe connects to a Windows named pipe.
func DialPipe(addr string, timeout time.Duration) (net.Conn, error) {
	return winio.DialPipe(addr, &timeout)
}
