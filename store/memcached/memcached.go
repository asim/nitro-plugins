package memcached

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	mc "github.com/bradfitz/gomemcache/memcache"
	"github.com/micro/go-micro/config/options"
	"github.com/micro/go-micro/store"
)

type mkv struct {
	options.Options
	Server *mc.ServerList
	Client *mc.Client
}

func (m *mkv) Read(key string) (*store.Record, error) {
	keyval, err := m.Client.Get(key)
	if err != nil && err == mc.ErrCacheMiss {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	if keyval == nil {
		return nil, store.ErrNotFound
	}

	return &store.Record{
		Key:    keyval.Key,
		Value:  keyval.Value,
		Expiry: time.Second * time.Duration(keyval.Expiration),
	}, nil
}

func (m *mkv) Delete(keys ...string) error {
	var err error
	for _, key := range keys {
		if err = m.Client.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

func (m *mkv) Write(record *store.Record) error {
	return m.Client.Set(&mc.Item{
		Key:        record.Key,
		Value:      record.Value,
		Expiration: int32(record.Expiry.Seconds()),
	})
}

func (m *mkv) Sync() ([]*store.Record, error) {
	// stats
	// cachedump
	// get keys

	var keys []string

	//store := make(map[string]string)
	if err := m.Server.Each(func(c net.Addr) error {
		cc, err := net.Dial("tcp", c.String())
		if err != nil {
			return err
		}
		defer cc.Close()

		b := bufio.NewReadWriter(bufio.NewReader(cc), bufio.NewWriter(cc))

		// get records
		if _, err := fmt.Fprintf(b, "stats records\r\n"); err != nil {
			return err
		}

		b.Flush()

		v, err := b.ReadSlice('\n')
		if err != nil {
			return err
		}

		parts := bytes.Split(v, []byte("\n"))
		if len(parts) < 1 {
			return nil
		}
		vals := strings.Split(string(parts[0]), ":")
		records := vals[1]

		// drain
		for {
			buf, err := b.ReadSlice('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if strings.HasPrefix(string(buf), "END") {
				break
			}
		}

		b.Writer.Reset(cc)
		b.Reader.Reset(cc)

		if _, err := fmt.Fprintf(b, "lru_crawler metadump %s\r\n", records); err != nil {
			return err
		}
		b.Flush()

		for {
			v, err := b.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if strings.HasPrefix(v, "END") {
				break
			}
			key := strings.Split(v, " ")[0]
			keys = append(keys, strings.TrimPrefix(key, "key="))
		}

		return nil
	}); err != nil {
		return nil, err
	}

	var vals []*store.Record

	// concurrent op
	ch := make(chan *store.Record, len(keys))

	for _, k := range keys {
		go func(key string) {
			i, _ := m.Read(key)
			ch <- i
		}(k)
	}

	for i := 0; i < len(keys); i++ {
		record := <-ch

		if record == nil {
			continue
		}

		vals = append(vals, record)
	}

	close(ch)

	return vals, nil
}

func (m *mkv) String() string {
	return "memcached"
}

func NewStore(opts ...options.Option) store.Store {
	options := options.NewOptions(opts...)

	var nodes []string

	if n, ok := options.Values().Get("store.nodes"); ok {
		nodes = n.([]string)
	}

	if len(nodes) == 0 {
		nodes = []string{"127.0.0.1:11211"}
	}

	ss := new(mc.ServerList)
	ss.SetServers(nodes...)

	return &mkv{
		Options: options,
		Server:  ss,
		Client:  mc.New(nodes...),
	}
}
