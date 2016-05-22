package kubernetes

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/micro/go-plugins/registry/kubernetes/client"

	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
)

type kregistry struct {
	sync.Mutex
	client  client.Kubernetes
	timeout time.Duration
}

type podServiceMeta struct {
	Service  *registry.Service
	Metadata map[string]string
}

var (
	// used within labels
	serviceLabelPrefix = "micro.mu/"
	serviceLabelKey    = "micro"
	// used within annotations
	serviceListKey = "micro/services"
	servicePrefix  = "micro/services/"
)

func init() {
	cmd.DefaultRegistries["kubernetes"] = NewRegistry
}

// Deregister will remove a known service from the Pod
func (c *kregistry) Deregister(s *registry.Service) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
	}

	// since we're getting the pod then updating, lets lock to serialise the update
	c.Lock()
	defer c.Unlock()

	// TODO: grab podname from somewhere better than this.
	podName := os.Getenv("HOSTNAME")

	// get the pod so we can update the existing list of services
	pod, err := c.client.GetPod(podName)
	if err != nil {
		return err
	}

	serviceData := pod.Metadata.Annotations[servicePrefix+s.Name]
	var services []*registry.Service

	if err := json.Unmarshal([]byte(serviceData), &services); err == nil {
		for i, service := range services {
			// only operate on the same version
			if service.Version != s.Version {
				continue
			}

			var nodes []*registry.Node

			for _, node := range service.Nodes {
				var seen bool
				for _, n := range s.Nodes {
					if n.Id == node.Id {
						seen = true
						break
					}
				}
				if !seen {
					nodes = append(nodes, node)
				}
			}

			// save
			service.Nodes = nodes
			services[i] = service
		}
	}

	// encode service
	b, err := json.Marshal(services)
	if err != nil {
		return err
	}

	pod = &client.Pod{
		Metadata: client.Meta{
			Annotations: map[string]string{
				servicePrefix + s.Name: string(b),
			},
		},
	}

	return c.client.UpdatePod(podName, pod)
}

// Register sets annotations on the current pod
func (c *kregistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
	}

	// since we're getting the pod then updating, lets lock to serialise the update
	c.Lock()
	defer c.Unlock()

	// TODO: grab podname from somewhere better than this.
	podName := os.Getenv("HOSTNAME")

	// get the pod so we can update the existing list of services
	pod, err := c.client.GetPod(podName)
	if err != nil {
		return err
	}

	// extract the service
	serviceData := pod.Metadata.Annotations[servicePrefix+s.Name]
	var services []*registry.Service

	if err := json.Unmarshal([]byte(serviceData), &services); err == nil {
		for i, service := range services {
			// only operate on the same version
			if service.Version != s.Version {
				continue
			}

			// iterate through old nodes to check if we're working on one
			for _, node := range service.Nodes {
				var seen bool

				// iterate through new nodes and check if
				// it exists, if so, mark and break
				for _, n := range s.Nodes {
					if n.Id == node.Id {
						seen = true
						break
					}
				}

				// if the node was not seen in the new list add it
				if !seen {
					s.Nodes = append(s.Nodes, node)
				}
			}

			// save
			services[i] = s
		}
	}

	// encode service
	b, err := json.Marshal(services)
	if err != nil {
		return err
	}

	pod = &client.Pod{
		Metadata: client.Meta{
			Labels: map[string]string{
				serviceLabelPrefix + s.Name: "service",
				serviceLabelKey:             "service",
			},
			Annotations: map[string]string{
				servicePrefix + s.Name: string(b),
			},
		},
	}

	return c.client.UpdatePod(podName, pod)
}

// GetService uses the `/endpoints/{name}` and `/pods?labelSelector=micro.mu/<name>=service`
// api to build the tree of Services, nodes, metadata and endpoints.
func (c *kregistry) GetService(name string) ([]*registry.Service, error) {
	podList, err := c.client.GetPods(map[string]string{serviceLabelPrefix + name: "service"})
	if err != nil {
		return nil, err
	}

	if len(podList.Items) == 0 {
		return nil, errors.New("No pods found with '" + serviceLabelPrefix + name + ":service' label")
	}

	svcs := make(map[string]*registry.Service)

	// loop through all pods and add to version map
	for _, pod := range podList.Items {
		v, ok := pod.Metadata.Annotations[servicePrefix+name]
		if !ok {
			continue
		}

		var srvs []*registry.Service

		if err := json.Unmarshal([]byte(v), &srvs); err != nil {
			continue
		}

		for _, srv := range srvs {
			service, ok := svcs[srv.Version]
			if !ok {
				svcs[srv.Version] = srv
				continue
			}
			// append nodes
			service.Nodes = append(service.Nodes, srv.Nodes...)
			svcs[srv.Version] = srv
		}
	}

	var services []*registry.Service

	for _, service := range svcs {
		services = append(services, service)
	}

	return services, nil
}

// ListServices will list all the service names
func (c *kregistry) ListServices() ([]*registry.Service, error) {
	// get all the pods labels as a micro service
	podList, err := c.client.GetPods(map[string]string{serviceLabelKey: "service"})
	if err != nil {
		return nil, err
	}

	if len(podList.Items) == 0 {
		return nil, errors.New("No pods found with '" + serviceLabelKey + ":service' label")
	}

	serviceMap := make(map[string]bool)

	// loop through all pods and get the services they know about
	for _, pod := range podList.Items {
		for k, _ := range pod.Metadata.Labels {
			if !strings.HasPrefix(k, serviceLabelPrefix) {
				continue
			}
			serviceMap[strings.TrimPrefix(k, serviceLabelPrefix)] = true
		}
	}

	var services []*registry.Service

	for service, _ := range serviceMap {
		services = append(services, &registry.Service{Name: service})
	}

	return services, nil
}

// Watch returns a kubernetes watcher
func (c *kregistry) Watch() (registry.Watcher, error) {
	return newWatcher(c)
}

func (c *kregistry) String() string {
	return "kubernetes"
}

// NewRegistry creates a kubernetes registry
func NewRegistry(opts ...registry.Option) registry.Registry {

	var options registry.Options
	for _, o := range opts {
		o(&options)
	}

	// get first host
	var host string
	if len(options.Addrs) > 0 && len(options.Addrs[0]) > 0 {
		host = options.Addrs[0]
	}

	if options.Timeout == 0 {
		options.Timeout = time.Second * 1
	}

	// if no hosts setup, assume InCluster
	var c client.Kubernetes
	if len(host) == 0 {
		c = client.NewClientInCluster()
	} else {
		c = client.NewClientByHost(host)
	}

	return &kregistry{
		client:  c,
		timeout: options.Timeout,
	}
}
