package nacos

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"

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

func getDeregisterTTL(t time.Duration) time.Duration {
	return 0
}

func newTransport(config *tls.Config) *http.Transport {
	return nil
}

func buildClientConfig(clientConfig *constant.ClientConfig, context context.Context) {
	if context != nil {
		if listenInterval, ok := context.Value("listen_interval").(uint64); ok {
			clientConfig.ListenInterval = listenInterval
		}
		if beatInterval, ok := context.Value("beat_interval").(int64); ok {
			clientConfig.BeatInterval = beatInterval
		}
		if namespaceId, ok := context.Value("namespace_id").(string); ok {
			clientConfig.NamespaceId = namespaceId
		}
		if endpoint, ok := context.Value("endpoint").(string); ok {
			clientConfig.Endpoint = endpoint
		}
		if accessKey, ok := context.Value("access_key").(string); ok {
			clientConfig.AccessKey = accessKey
		}
		if secretKey, ok := context.Value("secret_key").(string); ok {
			clientConfig.SecretKey = secretKey
		}
		if cacheDir, ok := context.Value("cache_dir").(string); ok {
			clientConfig.CacheDir = cacheDir
		}
		if logDir, ok := context.Value("log_dir").(string); ok {
			clientConfig.LogDir = logDir
		}
		if updateThreadNum, ok := context.Value("update_thread_num").(int); ok {
			clientConfig.UpdateThreadNum = updateThreadNum
		}
		if notLoadCacheAtStart, ok := context.Value("not_load_cache_at_start").(bool); ok {
			clientConfig.NotLoadCacheAtStart = notLoadCacheAtStart
		}
		if updateCacheWhenEmpty, ok := context.Value("update_cache_when_empty").(bool); ok {
			clientConfig.UpdateCacheWhenEmpty = updateCacheWhenEmpty
		}
		if openKMS, ok := context.Value("open_kms").(bool); ok {
			clientConfig.OpenKMS = openKMS
		}
		if regionId, ok := context.Value("region_id").(string); ok {
			clientConfig.RegionId = regionId
		}
		if username, ok := context.Value("user_name").(string); ok {
			clientConfig.Username = username
		}
		if password, ok := context.Value("password").(string); ok {
			clientConfig.Password = password
		}
	}
}

func configure(c *nacosRegistry, opts ...registry.Option) error {
	// set opts
	for _, o := range opts {
		o(&c.opts)
	}

	contentPath := "/nacos"
	if c.opts.Context != nil {
		if contextPath, ok := c.opts.Context.Value("context_path").(string); ok {
			contentPath = contextPath
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
	clientConfig := constant.ClientConfig{}
	clientConfig.TimeoutMs = uint64(c.opts.Timeout.Milliseconds())
	buildClientConfig(&clientConfig, c.opts.Context)
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
	if len(s.Nodes) == 0 {
		return errors.New("you must deregister at least one node")
	}
	var options registry.DeregisterOptions
	for _, o := range opts {
		o(&options)
	}
	node := s.Nodes[0]
	host, pt, err := net.SplitHostPort(node.Address)
	if err != nil {
		return err
	}
	port, err := strconv.Atoi(pt)
	if err != nil {
		return err
	}
	instance := vo.DeregisterInstanceParam{
		Ip:          host,
		Port:        uint64(port),
		ServiceName: s.Name,
	}
	if options.Context != nil {
		if tenant, ok := options.Context.Value("tenant").(string); ok {
			instance.Tenant = tenant
		}
		if cluster, ok := options.Context.Value("cluster").(string); ok {
			instance.Cluster = cluster
		}
		if groupName, ok := options.Context.Value("group_name").(string); ok {
			instance.GroupName = groupName
		}
		if ephemeral, ok := options.Context.Value("ephemeral").(bool); ok {
			instance.Ephemeral = ephemeral
		}
	}

	c.namingClient.DeregisterInstance(instance)
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
		Metadata:    s.Metadata,
		ServiceName: s.Name,
	}
	if options.Context != nil {
		if enable, ok := options.Context.Value("enable").(bool); ok {
			instance.Enable = enable
		}
		if healthy, ok := options.Context.Value("healthy").(bool); ok {
			instance.Healthy = healthy
		}
		if groupName, ok := options.Context.Value("group_name").(string); ok {
			instance.GroupName = groupName
		}
		if ephemeral, ok := options.Context.Value("ephemeral").(bool); ok {
			instance.Ephemeral = ephemeral
		}
		if tenant, ok := options.Context.Value("tenant").(string); ok {
			instance.Tenant = tenant
		}
		if cluster, ok := options.Context.Value("cluster").(string); ok {
			instance.ClusterName = cluster
		}
		if groupName, ok := options.Context.Value("group_name").(string); ok {
			instance.GroupName = groupName
		}
		if ephemeral, ok := options.Context.Value("ephemeral").(bool); ok {
			instance.Ephemeral = ephemeral
		}
		if weight, ok := options.Context.Value("weight").(float64); ok {
			instance.Weight = weight
		}
	}
	_, err := c.namingClient.RegisterInstance(instance)
	return err
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
	return c.opts
}

func NewRegistry(opts ...registry.Option) registry.Registry {
	nacos := &nacosRegistry{
		opts: registry.Options{},
	}
	configure(nacos, opts...)
	return nacos
}
