package grpc

import (
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

const (
	maxConcurrentStreamsChannel = 100
	cleanerSleep                = time.Second
)

type pool struct {
	sync.Mutex
	ttl     time.Duration
	running uint32
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
		ttl:   ttl,
		conns: make(map[string][]*poolConn),
	}

	return out
}

func (p *pool) cleaner() {
	for {
		time.Sleep(cleanerSleep)
		p.clearConns()

		if len(p.conns) == 0 {
			if atomic.CompareAndSwapUint32(&p.running, 1, 0) {
				break
			}
		}
	}

}

func (p *pool) wakeCleaner() {
	for atomic.LoadUint32(&p.running) != 1 {
		if atomic.CompareAndSwapUint32(&p.running, 0, 1) {
			go p.cleaner()
		}
	}
}

func (p *pool) clearConns() {
	p.Lock()
	defer p.Unlock()

	now := time.Now()
	copyConns := make(map[string][]*poolConn, len(p.conns))

	// the table save addr -> conns
	for addr, conns := range p.conns {
		copyCons := make([]*poolConn, 0, len(conns))

		// the array save conns
		for _, c := range conns {
			// don't release the connection
			if c.refCount > 0 || now.Sub(c.lastRefTime) < p.ttl {
				copyCons = append(copyCons, c)
				continue
			}

			// release useless connection
			c.client.Close()
		}

		// if the connections is not empty
		if len(copyCons) > 0 {
			copyConns[addr] = copyCons
		}
	}

	p.conns = copyConns
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

	p.wakeCleaner()

	return conn, nil
}

func (p *pool) release(con *poolConn) {
	p.Lock()
	defer p.Unlock()

	con.delRef()
}
