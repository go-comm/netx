package netx

import (
	"net"
	"time"
)

type timeoutConn struct {
	net.Conn
	options
}

func NewTimeoutConn(c net.Conn, opts ...Option) net.Conn {
	tc := &timeoutConn{Conn: c}
	for _, o := range opts {
		o(&tc.options)
	}
	return tc
}

func (c *timeoutConn) Read(b []byte) (n int, err error) {
	if c.ReadTimeout > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))
	}
	n, err = c.Conn.Read(b)
	return
}

func (c *timeoutConn) Write(b []byte) (n int, err error) {
	if c.WriteTimeout > 0 {
		c.Conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	}
	n, err = c.Conn.Write(b)
	return
}
