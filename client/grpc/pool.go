package grpc

import (
	"sync"
	"time"

	"google.golang.org/grpc"
)

type pool struct {
	sync.Mutex
	conns map[string]*poolConn
}

type poolConn struct {
	refCount    int32
	lastRefTime time.Time
	client      *grpc.ClientConn
}

func (c *poolConn) addRef() {
	c.lastRefTime = time.Now()
	c.refCount++
}

func (c *poolConn) delRef() {
	c.lastRefTime = time.Now()
	c.refCount--
}

func newPool() *pool {
	out := &pool{
		conns: make(map[string]*poolConn),
	}

	go func() {
		for {
			time.Sleep(time.Second)
			out.clear()
		}
	}()

	return out
}

func (p *pool) clear() {
	p.Lock()
	defer p.Unlock()

	now := time.Now()

	for k, v := range p.conns {
		if v.refCount > 0 {
			continue
		}

		if now.Sub(v.lastRefTime) < time.Minute*5 {
			continue
		}

		delete(p.conns, k)
		v.client.Close()
	}
}

func (p *pool) getConn(addr string, opts ...grpc.DialOption) (*poolConn, error) {
	p.Lock()
	defer p.Unlock()

	con, ok := p.conns[addr]
	if ok {
		con.addRef()
		return con, nil
	}

	// create new conn
	c, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	con = &poolConn{client: c}
	p.conns[addr] = con

	con.addRef()
	return con, nil
}

func (p *pool) release(con *poolConn) {
	p.Lock()
	defer p.Unlock()

	con.delRef()
}
