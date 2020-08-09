package nacos

import (
	reflect "reflect"
	"sync"

	mnet "github.com/micro/go-micro/v2/util/net"

	"github.com/micro/go-micro/v2/logger"

	"github.com/nacos-group/nacos-sdk-go/model"

	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/micro/go-micro/v2/registry"
)

type nacosWatcher struct {
	nr *nacosRegistry
	wo registry.WatchOptions

	next chan *registry.Result
	exit chan bool

	sync.RWMutex
	services      map[string][]*registry.Service
	cacheServices map[string][]model.SubscribeService
	param         *vo.SubscribeParam
}

func NewNacosWatcher(nr *nacosRegistry, opts ...registry.WatchOption) (registry.Watcher, error) {
	var wo registry.WatchOptions
	for _, o := range opts {
		o(&wo)
	}
	nw := nacosWatcher{
		nr:            nr,
		wo:            wo,
		exit:          make(chan bool),
		next:          make(chan *registry.Result, 10),
		services:      make(map[string][]*registry.Service),
		cacheServices: make(map[string][]model.SubscribeService),
		param:         new(vo.SubscribeParam),
	}
	if wo.Context != nil {
		if p, ok := wo.Context.Value("subscribe_param").(vo.SubscribeParam); ok {
			nw.param = &p
		}
	}
	nw.param.SubscribeCallback = nw.callBackHandle
	return &nw, nil
}

func (nw *nacosWatcher) callBackHandle(services []model.SubscribeService, err error) {
	if err != nil {
		logger.Error("nacos watcher call back handle error:%v", err)
		return
	}
	serviceName := services[0].ServiceName

	if nw.cacheServices[serviceName] == nil {
		for _, v := range services {
			nw.next <- &registry.Result{Action: "create", Service: buildRegistryService(&v)}
		}
	} else {
		for _, subscribeService := range services {
			create := true
			for _, cacheService := range nw.cacheServices[serviceName] {
				if subscribeService.InstanceId == cacheService.InstanceId {
					if !reflect.DeepEqual(subscribeService, cacheService) {
						nw.next <- &registry.Result{Action: "update", Service: buildRegistryService(&subscribeService)}
					}
					create = false
				}
			}
			if create {
				nw.next <- &registry.Result{Action: "create", Service: buildRegistryService(&subscribeService)}
			}
		}

		for _, cacheService := range nw.cacheServices[serviceName] {
			del := true
			for _, subscribeService := range services {
				if subscribeService.InstanceId == cacheService.InstanceId {
					del = false
				}
			}
			if del {
				nw.next <- &registry.Result{Action: "delete", Service: buildRegistryService(&cacheService)}
			}
		}
	}
	nw.cacheServices[serviceName] = services
}

func buildRegistryService(v *model.SubscribeService) (s *registry.Service) {
	nodes := make([]*registry.Node, 0)
	nodes = append(nodes, &registry.Node{
		Id:       v.InstanceId,
		Address:  mnet.HostPort(v.Ip, v.Port),
		Metadata: v.Metadata,
	})
	s = &registry.Service{
		Name:     v.ServiceName,
		Version:  "latest",
		Metadata: v.Metadata,
		Nodes:    nodes,
	}
	return
}

func (nw *nacosWatcher) Next() (r *registry.Result, err error) {
	select {
	case <-nw.exit:
		return nil, registry.ErrWatcherStopped
	case r, ok := <-nw.next:
		if !ok {
			return nil, registry.ErrWatcherStopped
		}
		return r, nil
	}
	// NOTE: This is a dead code path: e.g. it will never be reached
	// as we return in all previous code paths never leading to this return
	return nil, registry.ErrWatcherStopped
}

func (nw *nacosWatcher) Stop() {
	nw.nr.namingClient.Unsubscribe(nw.param)
	select {
	case <-nw.exit:
		return
	default:
		close(nw.exit)
	}
}
