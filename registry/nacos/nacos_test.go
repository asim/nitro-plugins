package nacos

import (
	"context"
	"testing"

	"github.com/nacos-group/nacos-sdk-go/model"

	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"

	"github.com/nacos-group/nacos-sdk-go/common/constant"

	"github.com/nacos-group/nacos-sdk-go/clients"

	"github.com/micro/go-micro/v2/registry"
)

type nacosClientMock struct {
	mock.Mock
}

func mockNamingClient() {

}

func (n *nacosClientMock) RegisterInstance(param vo.RegisterInstanceParam) (bool, error) {
	ret := n.Called(param)
	return ret.Bool(0), ret.Error(1)
}

func (n *nacosClientMock) DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error) {
	ret := n.Called(param)
	return ret.Bool(0), ret.Error(1)
}

func (n *nacosClientMock) GetService(param vo.GetServiceParam) (model.Service, error) {
	ret := n.Called(param)
	hosts := make([]model.Instance, 0)
	hosts = append(hosts, model.Instance{
		InstanceId:  "1",
		Ip:          "127.0.0.1",
		Port:        8080,
		Weight:      1.0,
		Metadata:    map[string]string{"version": "v1"},
		ServiceName: param.ServiceName,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   false,
	})
	service := model.Service{
		Name:  param.ServiceName,
		Hosts: hosts,
	}
	return service, ret.Error(1)
}

func (n *nacosClientMock) SelectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error) {
	ret := n.Called(param)
	hosts := make([]model.Instance, 0)
	hosts = append(hosts, model.Instance{
		InstanceId:  "1",
		Ip:          "127.0.0.1",
		Port:        8080,
		Weight:      1.0,
		Metadata:    map[string]string{"version": "v1"},
		ServiceName: param.ServiceName,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   false,
	})
	return hosts, ret.Error(1)

}

func (n *nacosClientMock) SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error) {
	ret := n.Called(param)
	hosts := make([]model.Instance, 0)
	hosts = append(hosts, model.Instance{
		InstanceId:  "1",
		Ip:          "127.0.0.1",
		Port:        8080,
		Weight:      1.0,
		Metadata:    map[string]string{"version": "v1"},
		ServiceName: param.ServiceName,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   false,
	})
	return hosts, ret.Error(1)
}

func (n *nacosClientMock) SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (*model.Instance, error) {
	ret := n.Called(param)
	return &model.Instance{
		InstanceId:  "1",
		Ip:          "127.0.0.1",
		Port:        8080,
		Weight:      1.0,
		Metadata:    map[string]string{"version": "v1"},
		ServiceName: param.ServiceName,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   false,
	}, ret.Error(1)
}

func (n *nacosClientMock) Subscribe(param *vo.SubscribeParam) error {
	ret := n.Called(param)
	return ret.Error(0)
}

func (n *nacosClientMock) Unsubscribe(param *vo.SubscribeParam) error {
	ret := n.Called(param)
	return ret.Error(0)
}

func (n *nacosClientMock) GetAllServicesInfo(param vo.GetAllServiceInfoParam) (model.ServiceList, error) {
	ret := n.Called(param)
	doms := make([]string, 2)
	doms[0] = "demo-service"
	doms[1] = "demo-service1"
	return model.ServiceList{
		Count: 1,
		Doms:  doms,
	}, ret.Error(1)
}

func buildNamingClient() (client naming_client.INamingClient, err error) {
	sc := []constant.ServerConfig{
		{
			IpAddr: "console.nacos.io",
			Port:   80,
		},
	}

	cc := constant.ClientConfig{
		TimeoutMs: 5000,
	}
	client, err = clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	return
}

func TestNacosNewRegistry(t *testing.T) {
	t.Run("NewRegistryWithContext", func(t *testing.T) {
		client, err := buildNamingClient()
		assert.Nil(t, err)
		r := NewRegistry(func(options *registry.Options) {
			options.Context = context.WithValue(options.Context, "naming_client", client)
		})
		assert.NotNil(t, r)
	})
	t.Run("NewRegistryWithAddr", func(t *testing.T) {
		r := NewRegistry(func(options *registry.Options) {
			options.Addrs = append(options.Addrs, "192.168.23.178:8848")
		})
		assert.NotNil(t, r)
	})
}

func TestNacosRegistry(t *testing.T) {
	t.Run("NacosRegistry", func(t *testing.T) {
		nacosClientMock := new(nacosClientMock)
		nacosClientMock.On("RegisterInstance", mock.Anything).Return(true, nil)
		r := NewRegistry(func(options *registry.Options) {
			options.Context = context.WithValue(options.Context, "naming_client", nacosClientMock)
		})
		assert.NotNil(t, r)
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
		client, err := buildNamingClient()
		assert.Nil(t, err)
		r := NewRegistry(func(options *registry.Options) {
			options.Context = context.WithValue(options.Context, "naming_client", client)
		})
		assert.NotNil(t, r)
		service := &registry.Service{}
		param := vo.RegisterInstanceParam{
			Ip:       "127.0.0.1",
			Port:     8080,
			Weight:   1.0,
			Enable:   true,
			Healthy:  true,
			Metadata: map[string]string{"version": "v1"},
		}
		err = r.Register(service, func(options *registry.RegisterOptions) {
			options.Context = context.WithValue(options.Context, "register_instance_param", param)
		})
		assert.Nil(t, err)
	})
}

func TestNacosDeRegistry(t *testing.T) {
	t.Run("NacosDeRegistry", func(t *testing.T) {
		nacosClientMock := new(nacosClientMock)
		nacosClientMock.On("DeregisterInstance", mock.Anything).Return(true, nil)
		r := NewRegistry(func(options *registry.Options) {
			options.Context = context.WithValue(options.Context, "naming_client", nacosClientMock)
		})
		assert.NotNil(t, r)
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
		client, err := buildNamingClient()
		assert.Nil(t, err)
		r := NewRegistry(func(options *registry.Options) {
			options.Context = context.WithValue(options.Context, "naming_client", client)
		})
		assert.NotNil(t, r)
		service := &registry.Service{}
		param := vo.DeregisterInstanceParam{
			Ip:          "127.0.0.1",
			Port:        8080,
			ServiceName: "demo",
		}
		err = r.Deregister(service, func(options *registry.DeregisterOptions) {
			options.Context = context.WithValue(options.Context, "deregister_instance_param", param)
		})
		assert.Nil(t, err)
	})
}
