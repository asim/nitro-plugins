package dns

import (
	"fmt"
	"os"

	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
)

const (
	ENV_DNS_SELECTOR_DOMAIN_NAME = "DNS_SELECTOR_DOMAIN_NAME"
	ENV_DNS_SELECTOR_PORT_NUMBER = "DNS_SELECTOR_PORT_NUMBER"
	DEFAULT_PORT_NUMBER          = "8080"
)

type dnsSelector struct {
	addressSuffix string
	envDomainName string
	envPortNumber string
}

func init() {
	cmd.DefaultSelectors["dns"] = NewSelector
}

func (s *dnsSelector) Init(opts ...selector.Option) error {
	return nil
}

func (s *dnsSelector) Options() selector.Options {
	return selector.Options{}
}

func (s *dnsSelector) Select(service string, opts ...selector.SelectOption) (selector.Next, error) {
	node := &registry.Node{
		Id:      service,
		Address: fmt.Sprintf("%v%v", service, s.addressSuffix),
	}

	return func() (*registry.Node, error) {
		return node, nil
	}, nil
}

func (s *dnsSelector) Mark(service string, node *registry.Node, err error) {
	return
}

func (s *dnsSelector) Reset(service string) {
	return
}

func (s *dnsSelector) Close() error {
	return nil
}

func (s *dnsSelector) String() string {
	return "named"
}

func NewSelector(opts ...selector.Option) selector.Selector {

	// Build a new
	s := &dnsSelector{
		addressSuffix: "",
		envDomainName: os.Getenv(ENV_DNS_SELECTOR_DOMAIN_NAME),
		envPortNumber: os.Getenv(ENV_DNS_SELECTOR_PORT_NUMBER),
	}

	// Add the dns domain-name (if one was specified by an env-var):
	if s.envDomainName != "" {
		s.addressSuffix += fmt.Sprintf(".%v", s.envDomainName)
	}

	// Either add the default port-number, or override with one specified by an env-var:
	if s.envPortNumber == "" {
		s.addressSuffix += fmt.Sprintf(":%v", DEFAULT_PORT_NUMBER)
	} else {
		s.addressSuffix += fmt.Sprintf(":%v", s.envPortNumber)
	}

	return s
}
