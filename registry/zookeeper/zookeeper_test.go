package zookeeper

import (
	"testing"

	"github.com/micro/go-micro/v3/registry"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
)

func testZKRegistry(t *testing.T) {
	reg := NewRegistry(
		registry.Addrs("127.0.0.1:2181"),
		registry.Timeout(20),
	)

	service := &registry.Service{
		Name:    "Test_Name",
		Version: "1.0.0",
		Metadata: map[string]string{
			"Header1": "H1",
		},
		Nodes: []*registry.Node{{
			Id:      "ID-123",
			Address: "127.0.0.1:8899",
			Metadata: map[string]string{
				"Header1": "Node1",
			}},
		},
	}

	Convey("test register", t, FailureHalts, func(c C) {
		err := reg.Register(service)
		c.So(err, ShouldBeNil)

		Convey("test getService", FailureHalts, func(c C) {
			ss, err := reg.GetService(service.Name)
			if err != nil {
				c.Printf("test getService err: %s\n", err)
				c.So(err, ShouldBeNil)
			}
			c.So(ss, should.HaveLength, 1)
			c.Printf("getService with name: %s\n", ss[0].Name)
		})

		Convey("test listService", FailureHalts, func(c C) {
			ss, err := reg.ListServices()
			if err != nil {
				c.Printf("test listService err: %s\n", err)
				c.So(err, ShouldBeNil)
			}
			c.So(ss, should.HaveLength, 1)
			c.Printf("listService with name: %s\n", ss[0].Name)
		})

		Convey("test deregister", FailureHalts, func(c C) {
			err := reg.Deregister(service)
			if err != nil {
				c.Printf("test deregister err: %s\n", err)
				c.So(err, ShouldBeNil)
			}
			Convey("test getService after deregister", FailureHalts, func(c C) {
				ss, err := reg.GetService(service.Name)
				if err != nil {
					c.Printf("test getService after deregister err: %s\n", err)
					c.So(err, ShouldBeNil)
				}
				c.So(ss, should.HaveLength, 0)
			})
		})
	})
}
