package grpc

import (
	"sync"

	"google.golang.org/grpc"
)

type pool struct {
	sync.Mutex
	conns map[string]*poolConn
}

type poolConn struct {
	client *grpc.ClientConn
}

func newPool() *pool {
	return &pool{
		conns: make(map[string]*poolConn),
	}
}

func (p *pool) getConn(addr string, opts ...grpc.DialOption) (*poolConn, error) {
	p.Lock()
	defer p.Unlock()

	con, ok := p.conns[addr]
	if ok {
		return con, nil
	}

	// create new conn
	c, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	con = &poolConn{c}
	p.conns[addr] = con

	return con, nil
}
