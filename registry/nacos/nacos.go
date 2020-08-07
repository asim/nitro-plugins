package nacos

import (
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"

	mnet "github.com/micro/go-micro/v2/util/net"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"

	"github.com/micro/go-micro/v2/cmd"
	"github.com/micro/go-micro/v2/registry"
)

type nacosRegistry struct {
	namingClient naming_client.INamingClient
	opts         registry.Options
}

func init() {
	cmd.DefaultRegistries["nacos"] = NewRegistry
}

func getNodeIpPort(s *registry.Service) (host string, port int, err error) {
	if len(s.Nodes) == 0 {
		return "", 0, errors.New("you must deregister at least one node")
	}
	node := s.Nodes[0]
	host, pt, err := net.SplitHostPort(node.Address)
	if err != nil {
		return "", 0, err
	}
	port, err = strconv.Atoi(pt)
	if err != nil {
		return "", 0, err
	}
	return
}

func configure(c *nacosRegistry, opts ...registry.Option) error {
	// set opts
	for _, o := range opts {
		o(&c.opts)
	}
	clientConfig := constant.ClientConfig{}
	contentPath := "/nacos"
	if c.opts.Context != nil {
		if contextPath, ok := c.opts.Context.Value("context_path").(string); ok {
			contentPath = contextPath
		}
		if clientConfig, ok := c.opts.Context.Value("client_config").(constant.ClientConfig); ok {
			clientConfig = clientConfig
		}
	}

	serverConfigs := make([]constant.ServerConfig, 0)

	// iterate the options addresses
	for _, address := range c.opts.Addrs {
		// check we have a port
		addr, port, err := net.SplitHostPort(address)
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			serverConfigs = append(serverConfigs, constant.ServerConfig{
				IpAddr:      addr,
				Port:        8848,
				ContextPath: contentPath,
			})
		} else if err == nil {
			p, err := strconv.ParseUint(port, 10, 64)
			if err != nil {
				continue
			}
			serverConfigs = append(serverConfigs, constant.ServerConfig{
				IpAddr:      addr,
				Port:        p,
				ContextPath: contentPath,
			})
		}
	}

	if c.opts.Timeout == 0 {
		c.opts.Timeout = time.Second * 1
	}
	clientConfig.TimeoutMs = uint64(c.opts.Timeout.Milliseconds())
	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
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
	host, port, err := getNodeIpPort(s)
	if err != nil {
		return err
	}
	var options registry.DeregisterOptions
	for _, o := range opts {
		o(&options)
	}

	param := vo.DeregisterInstanceParam{}
	if options.Context != nil {
		if p, ok := options.Context.Value("deregister_instance_param").(vo.DeregisterInstanceParam); ok {
			param = p
		}
	}
	param.Ip = host
	param.Port = uint64(port)
	param.ServiceName = s.Name

	_, err = c.namingClient.DeregisterInstance(param)
	return err
}

func (c *nacosRegistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	host, port, err := getNodeIpPort(s)
	if err != nil {
		return err
	}
	var options registry.RegisterOptions
	for _, o := range opts {
		o(&options)
	}
	// use first node

	param := vo.RegisterInstanceParam{}
	if options.Context != nil {
		if p, ok := options.Context.Value("deregister_instance_param").(vo.RegisterInstanceParam); ok {
			param = p
		}
	}
	param.Ip = host
	param.Port = uint64(port)
	param.Metadata = s.Metadata
	param.ServiceName = s.Name

	_, err = c.namingClient.RegisterInstance(param)
	return err
}

func (c *nacosRegistry) GetService(name string, opts ...registry.GetOption) ([]*registry.Service, error) {
	var options registry.GetOptions
	for _, o := range opts {
		o(&options)
	}
	param := vo.SelectInstancesParam{}
	if options.Context != nil {
		if p, ok := options.Context.Value("select_instances_param").(vo.SelectInstancesParam); ok {
			param = p
		}
	}
	param.ServiceName = name
	instances, err := c.namingClient.SelectInstances(param)
	if err != nil {
		return nil, err
	}
	services := make([]*registry.Service, 0)
	for _, v := range instances {
		nodes := make([]*registry.Node, 0)
		nodes = append(nodes, &registry.Node{
			Id:       v.InstanceId,
			Address:  mnet.HostPort(v.Ip, v.Port),
			Metadata: v.Metadata,
		})
		s := registry.Service{
			Name:     v.ServiceName,
			Version:  v.Metadata["version"],
			Metadata: v.Metadata,
			Nodes:    nodes,
		}
		services = append(services, &s)
	}

	return services, nil
}

func (c *nacosRegistry) ListServices(opts ...registry.ListOption) ([]*registry.Service, error) {
	var options registry.ListOptions
	for _, o := range opts {
		o(&options)
	}
	param := vo.GetAllServiceInfoParam{}
	if options.Context != nil {
		if p, ok := options.Context.Value("get_all_service_info_param").(vo.GetAllServiceInfoParam); ok {
			param = p
		}
	}
	services, err := c.namingClient.GetAllServicesInfo(param)
	if err != nil {
		return nil, err
	}
	var registryServices []*registry.Service
	for _, v := range services {
		registryServices = append(registryServices, &registry.Service{Name: v.Name})
	}
	return registryServices, nil
}

func (c *nacosRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	return nil, nil
}

func (c *nacosRegistry) String() string {
	return "nacos"
}

func (c *nacosRegistry) Options() registry.Options {
	return c.opts
}

func NewRegistry(opts ...registry.Option) registry.Registry {
	nacos := &nacosRegistry{
		opts: registry.Options{},
	}
	configure(nacos, opts...)
	return nacos
}
