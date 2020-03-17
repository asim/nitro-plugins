// Package consul is a consul implementation of kv
package consul

import (
	"fmt"
	"net"

	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2/store"
)

type ckv struct {
	options store.Options
	client  *api.Client
}

func (c *ckv) Init(...store.Option) error {
	return nil
}

func (c *ckv) Options() store.Options {
	return c.options
}

func (c *ckv) Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	// TODO: implement read options
	records := make([]*store.Record, 0, 1)

	keyval, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}

	if keyval == nil {
		return nil, store.ErrNotFound
	}

	records = append(records, &store.Record{
		Key:   keyval.Key,
		Value: keyval.Value,
	})

	return records, nil
}

func (c *ckv) Delete(key string, opts ...store.DeleteOption) error {
	_, err := c.client.KV().Delete(key, nil)
	return err
}

func (c *ckv) Write(record *store.Record, opts ...store.WriteOption) error {
	_, err := c.client.KV().Put(&api.KVPair{
		Key:   record.Key,
		Value: record.Value,
	}, nil)
	return err
}

func (c *ckv) List(opts ...store.ListOption) ([]string, error) {
	keyval, _, err := c.client.KV().List("/", nil)
	if err != nil {
		return nil, err
	}
	if keyval == nil {
		return nil, store.ErrNotFound
	}
	var keys []string
	for _, keyv := range keyval {
		keys = append(keys, keyv.Key)
	}
	return keys, nil
}

func (c *ckv) String() string {
	return "consul"
}

func NewStore(opts ...store.Option) store.Store {
	var options store.Options
	for _, o := range opts {
		o(&options)
	}

	config := api.DefaultConfig()
	nodes := options.Nodes

	// set host
	if len(nodes) > 0 {
		addr, port, err := net.SplitHostPort(nodes[0])
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			port = "8500"
			config.Address = fmt.Sprintf("%s:%s", nodes[0], port)
		} else if err == nil {
			config.Address = fmt.Sprintf("%s:%s", addr, port)
		}
	}

	client, _ := api.NewClient(config)

	return &ckv{
		options: options,
		client:  client,
	}
}
