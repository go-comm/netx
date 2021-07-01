package netx

import (
	"net"
	"sync"
)

func NewOnceCloseConn(c net.Conn) net.Conn {
	return &onceCloseConn{Conn: c}
}

type onceCloseConn struct {
	net.Conn
	once     sync.Once
	closeErr error
}

func (oc *onceCloseConn) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseConn) close() { oc.closeErr = oc.Conn.Close() }

func NewOnceCloseListener(l net.Listener) net.Listener {
	return &onceCloseListener{Listener: l}
}

type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() { oc.closeErr = oc.Listener.Close() }
