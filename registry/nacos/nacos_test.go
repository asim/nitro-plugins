package nacos

import (
	"context"
	"testing"

	"github.com/nacos-group/nacos-sdk-go/model"

	nacosMock "github.com/micro/go-plugins/registry/nacos/v2/mock"
	"github.com/stretchr/testify/mock"

	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/stretchr/testify/assert"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"

	"github.com/micro/go-micro/v2/registry"
)

func getRegistry(nacosClientMock naming_client.INamingClient) registry.Registry {
	r := NewRegistry(func(options *registry.Options) {
		options.Context = context.WithValue(options.Context, "naming_client", nacosClientMock)
	})
	return r
}

//nacos registry
func TestNacosRegistry(t *testing.T) {

	nacosClientMock := new(nacosMock.NacosClientMock)
	nacosClientMock.On("RegisterInstance", mock.Anything).Return(true, nil)

	r := getRegistry(nacosClientMock)
	assert.NotNil(t, r)

	t.Run("NacosRegistry", func(t *testing.T) {
		node := &registry.Node{
			Id:       "1",
			Address:  "127.0.0.1:8080",
			Metadata: map[string]string{"test": "test"},
		}
		nodes := make([]*registry.Node, 0)
		nodes = append(nodes, node)
		service := &registry.Service{
			Name:    "demo",
			Version: "latest",
			Nodes:   nodes,
		}
		err := r.Register(service)
		assert.Nil(t, err)
	})

	t.Run("NacosRegistryWithContext", func(t *testing.T) {
		service := &registry.Service{}
		param := vo.RegisterInstanceParam{
			Ip:        "127.0.0.1",
			Port:      8080,
			Weight:    1.0,
			Enable:    true,
			Healthy:   true,
			Metadata:  map[string]string{"version": "v1"},
			Ephemeral: true,
		}
		err := r.Register(service, func(options *registry.RegisterOptions) {
			options.Context = context.WithValue(options.Context, "register_instance_param", param)
		})
		assert.Nil(t, err)
	})
}

//nacos deregistry
func TestNacosDeRegistry(t *testing.T) {

	nacosClientMock := new(nacosMock.NacosClientMock)
	nacosClientMock.On("DeregisterInstance", mock.Anything).Return(true, nil)

	r := getRegistry(nacosClientMock)
	assert.NotNil(t, r)

	t.Run("NacosDeRegistry", func(t *testing.T) {
		node := &registry.Node{
			Id:       "1",
			Address:  "127.0.0.1:8080",
			Metadata: map[string]string{"test": "test"},
		}
		nodes := make([]*registry.Node, 0)
		nodes = append(nodes, node)
		service := &registry.Service{
			Name:    "demo",
			Version: "latest",
			Nodes:   nodes,
		}
		err := r.Deregister(service)
		assert.Nil(t, err)
	})

	t.Run("NacosDeRegistryWithContext", func(t *testing.T) {
		service := &registry.Service{}
		param := vo.DeregisterInstanceParam{
			Ip:          "127.0.0.1",
			Port:        8080,
			ServiceName: "demo",
		}
		err := r.Deregister(service, func(options *registry.DeregisterOptions) {
			options.Context = context.WithValue(options.Context, "deregister_instance_param", param)
		})
		assert.Nil(t, err)
	})
}

//nacos deregistry
func TestNacosGetService(t *testing.T) {
	nacosClientMock := new(nacosMock.NacosClientMock)
	nacosClientMock.On("GetService", mock.Anything).Return(model.Service{}, nil)

	r := getRegistry(nacosClientMock)
	assert.NotNil(t, r)

	t.Run("NacosGetService", func(t *testing.T) {
		services, err := r.GetService("demo")
		assert.True(t, len(services) == 1 && services[0].Name == "demo")
		assert.Nil(t, err)
	})

	t.Run("NacosGetServiceWithContext", func(t *testing.T) {
		param := vo.GetServiceParam{
			ServiceName: "demo",
			GroupName:   "DEFAULT_GROUP",
		}
		services, err := r.GetService("", func(options *registry.GetOptions) {
			options.Context = context.WithValue(options.Context, "select_instances_param", param)
		})
		assert.True(t, len(services) == 1 && services[0].Name == "demo")
		assert.Nil(t, err)
	})
}

//nacos deregistry
func TestNacosListServices(t *testing.T) {

	nacosClientMock := new(nacosMock.NacosClientMock)
	nacosClientMock.On("GetAllServicesInfo", mock.Anything).Return(model.ServiceList{}, nil)

	r := getRegistry(nacosClientMock)
	assert.NotNil(t, r)

	t.Run("NacosListServices", func(t *testing.T) {
		services, err := r.ListServices()
		assert.True(t, len(services) == 2 && err == nil)
	})

	t.Run("NacosListServicesWithContext", func(t *testing.T) {
		param := vo.GetAllServiceInfoParam{
			PageNo:   1,
			PageSize: 10,
		}
		services, err := r.ListServices(func(options *registry.ListOptions) {
			options.Context = context.WithValue(options.Context, "get_all_service_info_param", param)
		})
		assert.True(t, len(services) == 2 && err == nil)
	})
}
