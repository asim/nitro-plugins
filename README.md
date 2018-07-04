# Plugins [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/saturn-sonic/go-plugins?status.svg)](https://godoc.org/github.com/saturn-sonic/go-plugins) [![Travis CI](https://travis-ci.org/micro/go-plugins.svg?branch=master)](https://travis-ci.org/micro/go-plugins) [![Go Report Card](https://goreportcard.com/badge/micro/go-plugins)](https://goreportcard.com/report/github.com/saturn-sonic/go-plugins)

A repository for micro plugins

## Overview

Micro tooling is built on a powerful pluggable architecture. Plugins can be swapped out with zero code changes.
This repository contains plugins for all micro related tools. Read on for further info.

Check out the [Micro on NATS](https://micro.mu/blog/2016/04/11/micro-on-nats.html) blog post to learn more about plugins.

Follow us on [Twitter](https://twitter.com/microhq) or join the [Slack](http://slack.micro.mu/) community.

## Getting Started

* [Contents](#contents)
* [Usage](#usage)
* [Build Pattern](#build-pattern)
* [Contributions](#contributions)

## Contents

Contents of this repository:

| Directory | Description                                          |
| --------- | ---------------------------------------------------- |
| Broker    | PubSub messaging; NATS, NSQ, RabbitMQ, Kafka         |
| Client    | RPC Clients; gRPC, HTTP                              |
| Codec     | Message Encoding; BSON, Mercury                      |
| Micro     | Micro Toolkit Plugins                                |
| Registry  | Service Discovery; Etcd, Gossip, NATS                |
| Selector  | Load balancing; Label, Cache, Static                 |
| Server    | RPC Servers; gRPC, HTTP                              |
| Transport | Bidirectional Streaming; NATS, RabbitMQ              |
| Wrappers  | Middleware; Circuit Breakers, Rate Limiting, Tracing |

## Usage

Plugins can be added to go-micro in the following ways. By doing so they'll be available to set via command line args or environment variables.

### Import Plugins

```go
import (
	"github.com/micro/go-micro/cmd"
	_ "github.com/saturn-sonic/go-plugins/broker/rabbitmq"
	_ "github.com/saturn-sonic/go-plugins/registry/kubernetes"
	_ "github.com/saturn-sonic/go-plugins/transport/nats"
)

func main() {
	// Parse CLI flags
	cmd.Init()
}
```

The same is achieved when calling `service.Init`

```go
import (
	"github.com/micro/go-micro"
	_ "github.com/saturn-sonic/go-plugins/broker/rabbitmq"
	_ "github.com/saturn-sonic/go-plugins/registry/kubernetes"
	_ "github.com/saturn-sonic/go-plugins/transport/nats"
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
	"github.com/saturn-sonic/go-plugins/registry/kubernetes"
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

An anti-pattern is modifying the `main.go` file to include plugins. Best practice recommendation is to include
plugins in a separate file and rebuild with it included. This allows for automation of building plugins and
clean separation of concerns.

Create file plugins.go

```go
package main

import (
	_ "github.com/saturn-sonic/go-plugins/broker/rabbitmq"
	_ "github.com/saturn-sonic/go-plugins/registry/kubernetes"
	_ "github.com/saturn-sonic/go-plugins/transport/nats"
)
```

Build with plugins.go

```shell
go build -o service main.go plugins.go
```

Run with plugins

```shell
MICRO_BROKER=rabbitmq \
MICRO_REGISTRY=kubernetes \
MICRO_TRANSPORT=nats \
service
```

## Contributions

A few contributions by others

| Feature                                                                                  | Description                              | Author                                         |
| ---------------------------------------------------------------------------------------- | ---------------------------------------- | ---------------------------------------------- |
| [Registry/Kubernetes](https://godoc.org/github.com/saturn-sonic/go-plugins/registry/kubernetes) | Service discovery via the Kubernetes API | [@nickjackson](https://github.com/nickjackson) |
| [Registry/Zookeeper](https://godoc.org/github.com/saturn-sonic/go-plugins/registry/zookeeper)   | Service discovery using Zookeeper        | [@HeavyHorst](https://github.com/HeavyHorst)   |
