package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client represents a k8s client
type Client struct {
	Config *Config
	client *http.Client
}

// GetEndpoints returns a list of endpoints by name
func (k *Client) GetEndpoints(serviceName string) (*Endpoints, error) {
	var data Endpoints
	url, err := k.buildURL("endpoints/" + serviceName)
	if err != nil {
		return nil, err
	}

	if err := k.get(url, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// GetPod returns a single pod by name
func (k *Client) GetPod(name string) (*Pod, error) {
	url, err := k.buildURL("pods/" + name)
	if err != nil {
		return nil, err
	}
	var pod Pod
	if err := k.get(url, &pod); err != nil {
		return nil, err
	}
	return &pod, nil
}

// UpdatePod issues a PATCH to "/pods/{name}"
func (k *Client) UpdatePod(name string, pod *Pod) error {
	url, err := k.buildURL("pods/" + name)
	if err != nil {
		return err
	}

	if err := k.patch(url, pod); err != nil {
		return err
	}
	return nil
}

// GetPods finds all pods with the given labelSelectors
func (k *Client) GetPods(labels map[string]string) (*PodList, error) {
	qs := url.Values{}
	for k, v := range labels {
		qs.Add("labelSelector", k+"="+v)
	}
	url, err := k.buildURL("pods?" + qs.Encode())
	if err != nil {
		return nil, err
	}
	var data PodList
	if err := k.get(url, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// GetServices returns a list of services
func (k *Client) GetServices() (*ServiceList, error) {
	url, err := k.buildURL("services")
	if err != nil {
		return nil, err
	}
	var data ServiceList
	if err := k.get(url, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// WatchEndpoints creates a watchRequest for the `/endpoints` API
func (k *Client) WatchEndpoints() (*WatchRequest, error) {
	url, err := k.buildURL("endpoints")
	if err != nil {
		return nil, err
	}
	return newWatchRequest(k, url)
}

// patch encodes he data and issues a "application/merge-patch+json" request
func (k *Client) patch(url string, data interface{}) error {

	// Encode json
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(data); err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, b)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/merge-patch+json")
	if len(k.Config.BearerToken) > 0 {
		req.Header.Add("Authorization", "Bearer "+k.Config.BearerToken)
	}

	res, err := k.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return nil
}

// get issues a GET request for a given url
func (k *Client) get(url string, data interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	if len(k.Config.BearerToken) > 0 {
		req.Header.Add("Authorization", "Bearer "+k.Config.BearerToken)
	}

	res, err := k.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(data)
}

// buildURL will build a url from the host, namespace and parsed in string
func (k *Client) buildURL(r string) (string, error) {
	h := k.Config.Host
	n := k.Config.Namespace
	u, err := url.Parse(fmt.Sprintf("%s/api/v1/namespaces/%s/%s", h, n, r))
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
