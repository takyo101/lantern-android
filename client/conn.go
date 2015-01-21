package client

import (
	"errors"
	"net"
)

var (
	ErrCouldNotCreateListener = errors.New(`Could not create new listener.`)
)

// Listener is a wrapper around net.TCPListener that attempts to provide a
// Stop() function.
type Listener struct {
	*net.TCPListener
}

// NewListener creates a wrapper around TCPListener
func NewListener(addr string) (wrap *Listener, err error) {
	var li net.Listener
	var tli *net.TCPListener

	var ok bool

	if li, err = net.Listen("tcp", addr); err != nil {
		return nil, err
	}

	if tli, ok = li.(*net.TCPListener); !ok {
		return nil, ErrCouldNotCreateListener
	}

	wrap = &Listener{
		TCPListener: tli,
	}

	return wrap, nil
}

// Accept returns the next connection to the listener.
func (li *Listener) Accept() (net.Conn, error) {
	return li.TCPListener.Accept()
}

// Stop is currently not implemented but should make the listener stop
// accepting new connections and then kill all active connections.
func (li *Listener) Stop() error {
	// TODO
	return nil
}
