// Package etcdv3 provides an etcd version 3 registry
package etcdv3

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"path"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/clientv3"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"

	hash "github.com/mitchellh/hashstructure"
)

var (
	prefix      = "/micro-registry"
	cachePrefix = "/micro-cache"
)

type etcdv3Registry struct {
	client  *clientv3.Client
	options registry.Options
	sync.Mutex
	register map[string]uint64
	leases   map[string]clientv3.LeaseID
}

func init() {
	cmd.DefaultRegistries["etcdv3"] = NewRegistry
}

func encode(s *registry.Service) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func decode(ds []byte) *registry.Service {
	var s *registry.Service
	json.Unmarshal(ds, &s)
	return s
}

func nodePath(s, id string) string {
	service := strings.Replace(s, "/", "-", -1)
	node := strings.Replace(id, "/", "-", -1)
	return path.Join(prefix, service, node)
}

func servicePath(s string) string {
	return path.Join(prefix, strings.Replace(s, "/", "-", -1))
}

func (e *etcdv3Registry) Deregister(s *registry.Service) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
	}

	e.Lock()
	// delete our hash of the service
	delete(e.register, s.Name)
	// delete our lease of the service
	delete(e.leases, s.Name)
	e.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	for _, node := range s.Nodes {
		// cache the value for the watcher
		resp, err := e.client.Get(ctx, nodePath(s.Name, node.Id))
		if err != nil {
			return err
		}

		for _, ev := range resp.Kvs {
			lrsp, err := e.client.Grant(ctx, 5)
			if err != nil {
				return err
			}
			_, err = e.client.Put(ctx, path.Join(cachePrefix, string(ev.Key)), string(ev.Value), clientv3.WithLease(lrsp.ID))
			if err != nil {
				return err
			}
		}

		_, err = e.client.Delete(ctx, nodePath(s.Name, node.Id))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *etcdv3Registry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
	}

	//refreshing lease if existing
	leaseID, ok := e.leases[s.Name]
	if ok {
		_, err := e.client.KeepAliveOnce(context.TODO(), leaseID)
		if err != nil {
			return err
		}
	}

	var options registry.RegisterOptions
	for _, o := range opts {
		o(&options)
	}

	// create hash of service; uint64
	h, err := hash.Hash(s, nil)
	if err != nil {
		return err
	}

	// get existing hash
	e.Lock()
	v, ok := e.register[s.Name]
	e.Unlock()

	// the service is unchanged, skip registering
	if ok && v == h {
		return nil
	}

	service := &registry.Service{
		Name:      s.Name,
		Version:   s.Version,
		Metadata:  s.Metadata,
		Endpoints: s.Endpoints,
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	var lgr *clientv3.LeaseGrantResponse
	if options.TTL.Seconds() > 0 {
		lgr, err = e.client.Grant(ctx, int64(options.TTL.Seconds()))
		if err != nil {
			return err
		}
	}

	for _, node := range s.Nodes {
		service.Nodes = []*registry.Node{node}
		if lgr != nil {
			_, err = e.client.Put(ctx, nodePath(service.Name, node.Id), encode(service), clientv3.WithLease(lgr.ID))
		} else {
			_, err = e.client.Put(ctx, nodePath(service.Name, node.Id), encode(service))
		}
		if err != nil {
			return err
		}
	}

	e.Lock()
	// save our hash of the service
	e.register[s.Name] = h
	// save our leaseID of the service
	if lgr != nil {
		e.leases[s.Name] = lgr.ID
	}
	e.Unlock()

	return nil
}

func (e *etcdv3Registry) GetService(name string) ([]*registry.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	rsp, err := e.client.Get(ctx, servicePath(name)+"/", clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	if len(rsp.Kvs) == 0 {
		return nil, registry.ErrNotFound
	}

	serviceMap := map[string]*registry.Service{}

	for _, n := range rsp.Kvs {
		if sn := decode(n.Value); sn != nil {
			s, ok := serviceMap[sn.Version]
			if !ok {
				s = &registry.Service{
					Name:      sn.Name,
					Version:   sn.Version,
					Metadata:  sn.Metadata,
					Endpoints: sn.Endpoints,
				}
				serviceMap[s.Version] = s
			}

			for _, node := range sn.Nodes {
				s.Nodes = append(s.Nodes, node)
			}
		}
	}

	var services []*registry.Service
	for _, service := range serviceMap {
		services = append(services, service)
	}
	return services, nil
}

func (e *etcdv3Registry) ListServices() ([]*registry.Service, error) {
	var services []*registry.Service
	nameSet := make(map[string]struct{})

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	rsp, err := e.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	if len(rsp.Kvs) == 0 {
		return []*registry.Service{}, nil
	}

	for _, n := range rsp.Kvs {
		if sn := decode(n.Value); sn != nil {
			nameSet[sn.Name] = struct{}{}
		}
	}
	for k := range nameSet {
		service := &registry.Service{}
		service.Name = k
		services = append(services, service)
	}

	return services, nil
}

func (e *etcdv3Registry) Watch() (registry.Watcher, error) {
	return newEtcdv3Watcher(e, e.options.Timeout)
}

func (e *etcdv3Registry) String() string {
	return "etcdv3"
}

func NewRegistry(opts ...registry.Option) registry.Registry {
	config := clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}

	var options registry.Options
	for _, o := range opts {
		o(&options)
	}

	if options.Timeout == 0 {
		options.Timeout = 5 * time.Second
	}

	if options.Secure || options.TLSConfig != nil {
		tlsConfig := options.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		config.TLS = tlsConfig
	}

	var cAddrs []string

	for _, addr := range options.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}

	// if we got addrs then we'll update
	if len(cAddrs) > 0 {
		config.Endpoints = cAddrs
	}

	cli, _ := clientv3.New(config)
	e := &etcdv3Registry{
		client:   cli,
		options:  options,
		register: make(map[string]uint64),
		leases:   make(map[string]clientv3.LeaseID),
	}

	return e
}
