# Index Plugin

The index plugin is a micro toolkit plugin which enables the HTTP index to be routed to a service or respond with static content.

## Usage

Register the plugin before building Micro

```
package main

import (
	"github.com/micro/micro/plugin"
	"github.com/micro/go-plugins/api/index"
)

func init() {
	plugin.Register(index.NewPlugin())
}
```

It can then be applied on the command line like so. This will route to the greeter service.

```
micro --index_service=greeter
```

### Route to Service

Specifying `--index_service=` flag will route to a particular service

Below routes to go.micro.api.greeter given the default API namespace of go.micro.api. In the web case it will route to go.micro.web.greeter.

```
micro --index_service=greeter
```

Alternatively specify the service when registering the plugin

```
func init() {
        plugin.Register(index.WithService("greeter"))
}
```

Note: You can specify just the service if using the "proxy" handler. Any other requires service and method e.g --index_service=greeter/say/hello

### Static Content

Instead of routing to a service you may want to serve static content

Do so in the following way

```
micro --index_status=200 --index_header=Content-Type:text/plain --index_body="hello world"
```

The same can be achieved when registering the plugin

```
func init() {
        plugin.Register(index.WithResponse(
		200,
		http.Header{"Content-Type": []string{"text/plain"}},
		[]byte(`hello world`),
	))
}
```

### Scoped to API

If you like to only apply the plugin for a specific component you can register it with that specifically. 
For example, below you'll see the plugin registered with the API.

```
package main

import (
	"github.com/micro/micro/api"
	"github.com/micro/go-plugins/api/index"
)

func init() {
	api.Register(index.NewPlugin())
}
```

Here's what the help displays when you do that.

```
$ go run main.go plugin.go api --help
NAME:
   main api - Run the micro API

USAGE:
   main api [command options] [arguments...]

OPTIONS:
   --address 		Set the api address e.g 0.0.0.0:8080 [$MICRO_API_ADDRESS]
   --handler 		Specify the request handler to be used for mapping HTTP requests to services; {api, proxy, rpc} [$MICRO_API_HANDLER]
   --namespace 		Set the namespace used by the API e.g. com.example.api [$MICRO_API_NAMESPACE]
   --cors 		Comma separated whitelist of allowed origins for CORS [$MICRO_API_CORS]
   --index_service 	Service name to route index to. Specified without namespace e.g greeter [$INDEX_SERVICE]
   --index_status "0"	HTTP status code for response [$INDEX_STATUS]
   --index_header 	Comma separated list of key-value pairs for response header [$INDEX_HEADER]
   --index_body 	Body of the response [$INDEX_BODY]

```

In this case the usage would be

```
micro api --index_service=greeter
```
