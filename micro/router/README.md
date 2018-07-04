# Router Plugin

The router plugin is a HTTP handler plugin for the Micro API which enables you to define routes via go-os/config. This is 
dynamic configuration that can then be leveraged via anything that implements the go-os/config interface e.g file, etcd, consul 
or the config service.

## Features

- Request Matching
- Weighted Routing
- Reverse Proxying
- Priority Rules
- Configurable via go-os/Config
- Pluggable via micro/plugins

## TODO

- Regex Matching Host/Path

## Usage

Register the plugin before building Micro

```go
package main

import (
	"github.com/micro/micro/plugin"
	"github.com/saturn-sonic/go-plugins/micro/router"
)

func init() {
	plugin.Register(router.NewRouter())
}
```

## Config

Configuring the router is done via a go-os/Config source. Here's an example using the File source.

```go
// Create Config Source
f := file.NewSource(
	// Use routes.json file
	config.SourceName("routes.json"),
)

// Create Config
c := config.NewConfig(
	// With Source
	config.WithSource(f),
)

// Create Router
r := router.NewRouter(
	// With Config
	router.Config(c),
)
```

## Routes

Routes are used to config request to match and the response to return. Here's an example.

```json
{
	"api": {
		"routes": [
			{
				"request": {
					"method": "GET",
					"host": "127.0.0.1:10001",
					"path": "/"
				},
				"response": {
					"status_code": 302,
					"header": {
						"location": "http://example.com"
					}
				},
				"weight": 1.0
			},
			{
				"request": {
					"method": "POST",
					"host": "127.0.0.1:10001",
					"path": "/foo"
				},
				"response": {
					"status_code": 301,
					"header": {
						"location": "http://foo.bar.com"
					}
				},
				"weight": 1.0
			}
		]
	}
}
```
