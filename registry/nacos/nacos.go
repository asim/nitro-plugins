package nacos

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"

	"github.com/nacos-group/nacos-sdk-go/clients/config_client"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"

	"github.com/micro/go-micro/v2/cmd"
	"github.com/micro/go-micro/v2/registry"
)

type nacosRegistry struct {
	namingClient naming_client.INamingClient
	configClient config_client.IConfigClient
	options      registry.Options
}

func init() {
	cmd.DefaultRegistries["nacos"] = NewRegistry
}

func getDeregisterTTL(t time.Duration) time.Duration {
	return 0
}

func newTransport(config *tls.Config) *http.Transport {
	return nil
}

func configure(c *nacosRegistry, opts ...registry.Option) error {
	// set opts
	for _, o := range opts {
		o(&c.options)
	}
	// get first host
	var host string
	if len(c.options.Addrs) > 0 && len(c.options.Addrs[0]) > 0 {
		host = c.options.Addrs[0]
	}

	if c.options.Timeout == 0 {
		c.options.Timeout = time.Second * 1
	}

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": []constant.ServerConfig{
			{
				IpAddr: host,
				Port:   8848,
			},
		},
		"clientConfig": constant.ClientConfig{
			TimeoutMs:           uint64(c.options.Timeout.Milliseconds()),
			ListenInterval:      10000,
			NotLoadCacheAtStart: true,
			LogDir:              "data/nacos/log",
		},
	})
	if err != nil {
		return err
	}
	c.namingClient = client
	return nil
}

func (c *nacosRegistry) Init(opts ...registry.Option) error {
	return configure(c, opts...)
}

func (c *nacosRegistry) Deregister(s *registry.Service, opts ...registry.DeregisterOption) error {
	return nil
}

func (c *nacosRegistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	if len(s.Nodes) == 0 {
		return errors.New("you must register at least one node")
	}
	var options registry.RegisterOptions
	for _, o := range opts {
		o(&options)
	}
	// use first node
	node := s.Nodes[0]
	host, pt, _ := net.SplitHostPort(node.Address)
	if host == "" {
		host = node.Address
	}
	port, _ := strconv.Atoi(pt)
	instance := vo.RegisterInstanceParam{
		Ip:          host,
		Port:        uint64(port),
		Enable:      true,
		Healthy:     true,
		Metadata:    s.Metadata,
		ServiceName: s.Name,
	}
	c.namingClient.RegisterInstance(instance)
	return nil
}

func (c *nacosRegistry) GetService(name string, opts ...registry.GetOption) ([]*registry.Service, error) {
	return nil, nil
}

func (c *nacosRegistry) ListServices(opts ...registry.ListOption) ([]*registry.Service, error) {
	return nil, nil
}

func (c *nacosRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	return nil, nil
}

func (c *nacosRegistry) String() string {
	return "nacos"
}

func (c *nacosRegistry) Options() registry.Options {
	return registry.Options{}
}

func NewRegistry(opts ...registry.Option) registry.Registry {
	nacos := &nacosRegistry{
		options: registry.Options{},
	}
	configure(nacos, opts...)
	return nacos
}
