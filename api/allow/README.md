# Allow Plugin

The allow plugin is a plugin for micro which allows the services that can be used via the /rpc HTTP endpoint.

## Usage

Register the plugin before building Micro

```
package main

import (
	"github.com/micro/micro/plugin"
	"github.com/micro/go-plugins/api/allow"
)

func init() {
	plugin.Register(allow.NewRPCAllow())
}
```

It can then be applied on the command line like so.

```
micro --rpc_allow go.micro.srv.greeter,go.micro.srv.example api
```

### Scoped to API

If you like to only apply the plugin for a specific component you can register it with that specifically. 
For example, below you'll see the plugin registered with the API.

```
package main

import (
	"github.com/micro/micro/api"
	"github.com/micro/go-plugins/api/allow"
)

func init() {
	api.Register(allow.NewRPCAllow())
}
```

Here's what the help displays when you do that.

```
$ go run main.go link.go api --help
NAME:
   main api - Run the micro API

USAGE:
   main api [command options] [arguments...]

OPTIONS:
   --rpc_allow 	Comma separated allow of allowed services for RPC calls [$MICRO_RPC_ALLOW]
```

In this case the usage would be

```
micro api --rpc_allow go.micro.srv.greeter
```
