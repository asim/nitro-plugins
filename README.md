# Go Plugins [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/micro/go-plugins?status.svg)](https://godoc.org/github.com/micro/go-plugins) [![Travis CI](https://travis-ci.org/micro/go-plugins.svg?branch=master)](https://travis-ci.org/micro/go-plugins) [![Go Report Card](https://goreportcard.com/badge/micro/go-plugins)](https://goreportcard.com/report/github.com/micro/go-plugins)

A repository for go-micro plugins.

Check out the [Micro on NATS](https://blog.micro.mu/2016/04/11/micro-on-nats.html) blog post to learn more about plugins.

Follow us on Twitter at [@MicroHQ](https://twitter.com/microhq), join the [Slack](https://micro-services.slack.com) community [here](http://slack.micro.mu/) or 
check out the [Mailing List](https://groups.google.com/forum/#!forum/microhq).

## What's here?

Directory	|	Description
---		|	---
Broker		|	Asynchronous Pub/Sub; NATS, NSQ, RabbitMQ, Kafka	
Client		|	Alternative clients; gRPC, HTTP
Codec		|	RPC Encoding; BSON, Mercury
KV		|	Key-Value; Memcached, Redis
Metrics		|	Instrumentation; Statsd, Telegraf, Prometheus
Micro		|	Micro Toolkit Plugins
Registry	|	Service Discovery; Etcd, Gossip, NATS
Selector	|	Node Selection; Label, Mercury
Server		|	Alternative servers; gRPC, HTTP
Sync		|	Locking/Leadership election; Consul, Etcd
Trace		|	Distributed tracing; Zipkin
Transport	|	Synchronous Request/Response; NATS, RabbitMQ
Wrappers	|	Client/Server middleware; Circuit Breakers, Rate Limit


## Community Contributions

Feature		|	Description		|	Author
----------	|	------------		|	--------
[Registry/Kubernetes](https://godoc.org/github.com/micro/go-plugins/registry/kubernetes)	|	Service discovery via the Kubernetes API	|	[@nickjackson](https://github.com/nickjackson)
[Registry/Zookeeper](https://godoc.org/github.com/micro/go-plugins/registry/zookeeper)	|	Service discovery using Zookeeper	|	[@HeavyHorst](https://github.com/HeavyHorst)

## Usage

Plugins can be added to go-micro in the following ways. By doing so they'll be available to set via command line args or environment variables.

### Import Plugins

```go
import (
	"github.com/micro/go-micro/cmd"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/transport/nats"
)

func main() {
	// Parse CLI flags
	cmd.Init()
}
```

The same is achieved when calling ```service.Init```

```go
import (
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/transport/nats"
)

func main() {
	service := micro.NewService(
		// Set service name
		micro.Name("my.service"),
	)

	// Parse CLI flags
	service.Init()
}
```

### Use via CLI Flags

Activate via a command line flag

```shell
go run service.go --broker=rabbitmq --registry=kubernetes --transport=nats
```

### Use Plugins Directly

CLI Flags provide a simple way to initialise plugins but you can do the same yourself.

```go
import (
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/registry/kubernetes"
)

func main() {
	registry := kubernetes.NewRegistry() //a default to using env vars for master API

	service := micro.NewService(
		// Set service name
		micro.Name("my.service"),
		// Set service registry
		micro.Registry(registry),
	)
}
```

## Build Pattern

You may want to swap out plugins using automation or add plugins to the micro toolkit. 
An easy way to do this is by maintaining a separate file for plugin imports and including it during the build.

Create file plugins.go
```go
package main

import (
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/transport/nats"
)
```

Build with plugins.go
```shell
go build -o service main.go plugins.go
```

Run with plugins
```shell
service --broker=rabbitmq --registry=kubernetes --transport=nats
```
