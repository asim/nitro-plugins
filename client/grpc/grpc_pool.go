package grpc

import (
	"sync"
	"time"

	"google.golang.org/grpc"
)

const (
	maxConcurrentStreamsChannel = 100
	cleanerSleep                = time.Second
)

type pool struct {
	sync.Mutex
	idleTTL time.Duration
	conns   map[string][]*poolConn
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

func newPool(ttl time.Duration) *pool {
	out := &pool{
		idleTTL: ttl,
		conns:   make(map[string][]*poolConn),
	}

	go func() {
		for {
			time.Sleep(cleanerSleep)
			out.clear()
		}
	}()

	return out
}

func (p *pool) clear() {
	p.Lock()
	defer p.Unlock()

	now := time.Now()

	for addr, conns := range p.conns {
		for idx, c := range conns {
			if c.refCount > 0 {
				continue
			}

			if now.Sub(c.lastRefTime) < p.idleTTL {
				continue
			}

			conns = append(conns[:idx], conns[idx+1:]...)
			p.conns[addr] = conns
			c.client.Close()
		}
	}
}

func (p *pool) getIdleConn(addr string) *poolConn {
	p.Lock()
	defer p.Unlock()

	conns := p.conns[addr]

	for _, conn := range conns {
		if conn.refCount < maxConcurrentStreamsChannel {
			conn.addRef()
			return conn
		}
	}

	return nil
}

func (p *pool) getConn(addr string, opts ...grpc.DialOption) (*poolConn, error) {
	conn := p.getIdleConn(addr)
	if conn != nil {
		return conn, nil
	}

	// create new conn
	c, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	conn = &poolConn{client: c}
	conn.addRef()

	p.Lock()
	defer p.Unlock()

	conns := p.conns[addr]
	conns = append(conns, conn)
	p.conns[addr] = conns

	return conn, nil
}

func (p *pool) release(con *poolConn) {
	p.Lock()
	defer p.Unlock()

	con.delRef()
}
