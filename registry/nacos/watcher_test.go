package nacos

import (
	"context"
	"testing"

	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/micro/go-micro/v2/registry"
	nacosMock "github.com/micro/go-plugins/registry/nacos/v2/mock"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

func TestNacosRegistryWatch(t *testing.T) {
	nacosClientMock := new(nacosMock.NacosClientMock)
	nacosClientMock.On("Subscribe", mock.Anything).Return(nil)
	nacosClientMock.On("GetAllServicesInfo", mock.Anything).Return(model.ServiceList{}, nil)

	r := getRegistry(nacosClientMock)
	assert.NotNil(t, r)
	//watch all service
	t.Run("NacosRegistryWatch", func(t *testing.T) {
		watcher, err := r.Watch()
		assert.NotNil(t, watcher)
		assert.Nil(t, err)
	})

	//watch single service
	t.Run("NacosRegistryWatchWithContext", func(t *testing.T) {
		param := vo.SubscribeParam{
			ServiceName: "demo",
		}
		watcher, err := r.Watch(func(options *registry.WatchOptions) {
			options.Context = context.WithValue(options.Context, "subscribe_param", param)
		})
		assert.NotNil(t, watcher)
		assert.Nil(t, err)
	})
}

func TestCallBackHandle(t *testing.T) {
	nw := nacosWatcher{
		nr:            &nacosRegistry{},
		exit:          make(chan bool),
		next:          make(chan *registry.Result, 10),
		services:      make(map[string][]*registry.Service),
		cacheServices: make(map[string][]model.SubscribeService),
	}
	//Create action
	t.Run("CallBackHandleCreate", func(t *testing.T) {
		services := make([]model.SubscribeService, 1)
		services[0] = model.SubscribeService{
			InstanceId:  "1",
			Ip:          "127.0.0.1",
			Port:        1234,
			ServiceName: "DEMO",
		}
		nw.callBackHandle(services, nil)
		result, err := nw.Next()
		assert.True(t, result.Action == "create" && result.Service.Name == "DEMO" && err == nil)
	})

	//Update action
	t.Run("CallBackHandleUpdate", func(t *testing.T) {
		services := make([]model.SubscribeService, 1)
		services[0] = model.SubscribeService{
			InstanceId:  "1",
			Ip:          "127.0.0.1",
			Port:        1234,
			ServiceName: "DEMO1",
		}
		nw.callBackHandle(services, nil)
		result, err := nw.Next()
		assert.True(t, result.Action == "create" && result.Service.Name == "DEMO1" && err == nil)
		services = make([]model.SubscribeService, 1)
		services[0] = model.SubscribeService{
			InstanceId:  "1",
			Ip:          "127.0.0.1",
			Port:        12345,
			ServiceName: "DEMO1",
		}
		nw.callBackHandle(services, nil)
		result, err = nw.Next()
		assert.True(t, result.Action == "update" && result.Service.Name == "DEMO1" && result.Service.Nodes[0].Address == "127.0.0.1:12345")
		assert.Nil(t, err)
	})

	//Delete action
	t.Run("CallBackHandleDelete", func(t *testing.T) {
		services := make([]model.SubscribeService, 2)
		services[0] = model.SubscribeService{
			InstanceId:  "1",
			Ip:          "127.0.0.1",
			Port:        1234,
			ServiceName: "DEMO1",
		}
		services[1] = model.SubscribeService{
			InstanceId:  "2",
			Ip:          "127.0.0.1",
			Port:        12345,
			ServiceName: "DEMO1",
		}
		nw.callBackHandle(services, nil)
		result, err := nw.Next()
		assert.True(t, result.Action == "create" && result.Service.Name == "DEMO1" && err == nil)
		services = make([]model.SubscribeService, 1)
		services[0] = model.SubscribeService{
			InstanceId:  "1",
			Ip:          "127.0.0.1",
			Port:        1234,
			ServiceName: "DEMO1",
		}
		nw.callBackHandle(services, nil)
		result, err = nw.Next()
		assert.True(t, result.Action == "delete" && result.Service.Name == "DEMO1" && result.Service.Nodes[0].Address == "127.0.0.1:12345")
		assert.Nil(t, err)
	})

}

func TestWatchStop(t *testing.T) {
	nacosClientMock := new(nacosMock.NacosClientMock)
	nacosClientMock.On("Unsubscribe", mock.Anything).Return(nil)
	doms := make([]string, 2)
	doms[0] = "DEMO"
	doms[1] = "DEMO1"
	nw := nacosWatcher{
		nr: &nacosRegistry{
			namingClient: nacosClientMock,
		},
		exit:          make(chan bool),
		next:          make(chan *registry.Result, 10),
		services:      make(map[string][]*registry.Service),
		cacheServices: make(map[string][]model.SubscribeService),
		Doms:          doms,
	}
	t.Run("WatchStop", func(t *testing.T) {
		nw.Stop()
		_, isOpen := <-nw.exit
		assert.True(t, isOpen == false)
	})
}
