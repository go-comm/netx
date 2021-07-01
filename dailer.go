package netx

import (
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/xerrors"
)

func NewDailConnect(dail func() (net.Conn, error), opts ...Option) net.Conn {
	c := &DailConnect{
		Dail: dail,
	}
	for _, o := range opts {
		o(&c.options)
	}
	return c
}

var _ net.Conn = (*DailConnect)(nil)

type DailConnect struct {
	mutex sync.RWMutex
	conn  net.Conn
	Dail  func() (net.Conn, error)
	options
}

func (c *DailConnect) connect(reconnect bool) (conn net.Conn, err error) {
	conn = c.conn
	if reconnect || conn == nil {
		c.mutex.Lock()
		conn = c.conn
		if reconnect || conn == nil {
			if reconnect && conn != nil {
				conn.Close()
			}
			conn, err := c.Dail()
			if err != nil {
				c.conn = nil
				c.mutex.Unlock()
				return nil, err
			}
			c.conn = conn
		}
		c.mutex.Unlock()
	}
	return conn, nil
}

func (c *DailConnect) Read(b []byte) (n int, err error) {
	var conn net.Conn
	conn, err = c.connect(false)
	if err != nil {
		return 0, err
	}
	if c.ReadTimeout > 0 {
		conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))
	}
	n, err = conn.Read(b)
	if err != nil && xerrors.Is(err, os.ErrDeadlineExceeded) {
		if c.Reconnect {
			c.connect(c.Reconnect)
			return 0, nil
		}
		return
	}
	return
}

func (c *DailConnect) Write(b []byte) (n int, err error) {
	var conn net.Conn
	conn, err = c.connect(false)
	if err != nil {
		return 0, err
	}
	if c.WriteTimeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	}
	n, err = conn.Write(b)
	if err != nil && xerrors.Is(err, os.ErrDeadlineExceeded) {
		if c.Reconnect {
			c.connect(c.Reconnect)
			return 0, nil
		}
		return
	}
	return
}

func (c *DailConnect) Close() error {
	conn, err := c.connect(false)
	if err != nil {
		return err
	}
	return conn.Close()
}

func (c *DailConnect) LocalAddr() net.Addr {
	conn, err := c.connect(false)
	if err != nil {
		return nil
	}
	return conn.LocalAddr()
}

func (c *DailConnect) RemoteAddr() net.Addr {
	conn, err := c.connect(false)
	if err != nil {
		return nil
	}
	return conn.RemoteAddr()
}

func (c *DailConnect) SetDeadline(t time.Time) error {
	conn, err := c.connect(false)
	if err != nil {
		return err
	}
	return conn.SetDeadline(t)
}

func (c *DailConnect) SetReadDeadline(t time.Time) error {
	conn, err := c.connect(false)
	if err != nil {
		return err
	}
	return conn.SetReadDeadline(t)
}

func (c *DailConnect) SetWriteDeadline(t time.Time) error {
	conn, err := c.connect(false)
	if err != nil {
		return err
	}
	return conn.SetWriteDeadline(t)
}
