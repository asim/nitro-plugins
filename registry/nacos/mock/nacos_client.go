package mock

import (
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/stretchr/testify/mock"
)

type NacosClientMock struct {
	mock.Mock
}

func (n *NacosClientMock) RegisterInstance(param vo.RegisterInstanceParam) (bool, error) {
	ret := n.Called(param)
	return ret.Bool(0), ret.Error(1)
}

func (n *NacosClientMock) DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error) {
	ret := n.Called(param)
	return ret.Bool(0), ret.Error(1)
}

func (n *NacosClientMock) GetService(param vo.GetServiceParam) (model.Service, error) {
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

func (n *NacosClientMock) SelectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error) {
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

func (n *NacosClientMock) SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error) {
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

func (n *NacosClientMock) SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (*model.Instance, error) {
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

func (n *NacosClientMock) Subscribe(param *vo.SubscribeParam) error {
	ret := n.Called(param)
	return ret.Error(0)
}

func (n *NacosClientMock) Unsubscribe(param *vo.SubscribeParam) error {
	ret := n.Called(param)
	return ret.Error(0)
}

func (n *NacosClientMock) GetAllServicesInfo(param vo.GetAllServiceInfoParam) (model.ServiceList, error) {
	ret := n.Called(param)
	doms := make([]string, 2)
	doms[0] = "demo-service"
	doms[1] = "demo-service1"
	return model.ServiceList{
		Count: 2,
		Doms:  doms,
	}, ret.Error(1)
}
