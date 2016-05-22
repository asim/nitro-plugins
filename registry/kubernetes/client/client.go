package client

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

// Kubernetes ...
type Kubernetes interface {
	GetEndpoints(serviceName string) (*Endpoints, error)
	UpdatePod(podName string, pod *Pod) error
	GetPod(name string) (*Pod, error)
	GetPods(labels map[string]string) (*PodList, error)
	GetServices() (*ServiceList, error)
	WatchEndpoints() (*WatchRequest, error)
}

// Config ...
type Config struct {
	Host        string
	Namespace   string
	BearerToken string
	// TLSClientConfig *tls.Config
	// Transport       *http.Transport
}

// NewClientByHost sets up a client by host
func NewClientByHost(host string) *Client {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{},
		DisableCompression: true,
	}

	c := &http.Client{
		Transport: tr,
	}

	return &Client{
		Config: &Config{
			Host:      host,
			Namespace: "default",
		},
		client: c,
	}
}

// NewClientInCluster should work similarily to the official api
// NewInClient by setting up a client configuration for use within
// a k8s pod.
func NewClientInCluster() *Client {
	host := "https://" + os.Getenv("KUBERNETES_SERVICE_HOST") + ":" + os.Getenv("KUBERNETES_SERVICE_PORT")
	sa := "/var/run/secrets/kubernetes.io/serviceaccount"

	s, err := os.Stat(sa)
	if err != nil {
		log.Fatal(err)
	}
	if s == nil || !s.IsDir() {
		log.Fatal(errors.New("no k8s service account found"))
	}

	token, err := ioutil.ReadFile(path.Join(sa, "token"))
	if err != nil {
		log.Fatal(err)
	}

	crt, err := CertPoolFromFile(path.Join(sa, "ca.crt"))
	if err != nil {
		log.Fatal(err)
	}

	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: crt,
			},
			DisableCompression: true,
		},
	}

	return &Client{
		Config: &Config{
			Host:        host,
			Namespace:   "default",
			BearerToken: string(token),
		},
		client: c,
	}
}

// ServiceList is the top level item for the
type ServiceList struct {
	Items []Item `json:"items"`
}

// Item ...
type Item struct {
	Metadata Meta `json:"metadata"`
}

// Meta ...
type Meta struct {
	Name        string            `json:"name,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Endpoints is the top level item for the /endpoints API
type Endpoints struct {
	Metadata Meta     `json:"metadata"`
	Subsets  []Subset `json:"subsets"`
}

// PodList ...
type PodList struct {
	Items []Pod `json:"items"`
}

// Pod is the top level item for a pod
type Pod struct {
	Metadata Meta `json:"metadata"`
}

// Subset ...
type Subset struct {
	Addresses []Address `json:"addresses"`
	Ports     []Port    `json:"ports"`
}

// Address ...
type Address struct {
	TargetRef ObjectRef `json:"targetRef"`
	IP        string    `json:"ip"`
}

// ObjectRef ... used as part of the endpoint address to identify pod
type ObjectRef struct {
	Name string
}

// Port ...
type Port struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

// Event ...
type Event struct {
	Type   string    `json:"type"`
	Object Endpoints `json:"object"`
}
