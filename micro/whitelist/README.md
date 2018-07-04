# Whitelist Plugin

The whitelist plugin is a plugin for micro which whitelists the services that can be used via the /rpc HTTP endpoint.

## Usage

Register the plugin before building Micro

```
package main

import (
	"github.com/micro/micro/plugin"
	"github.com/saturn-sonic/go-plugins/micro/whitelist"
)

func init() {
	plugin.Register(whitelist.NewRPCWhitelist())
}
```

It can then be applied on the command line like so.

```
micro --rpc_whitelist go.micro.srv.greeter,go.micro.srv.example api
```

### Scoped to API

If you like to only apply the plugin for a specific component you can register it with that specifically. 
For example, below you'll see the plugin registered with the API.

```
package main

import (
	"github.com/micro/micro/api"
	"github.com/saturn-sonic/go-plugins/micro/whitelist"
)

func init() {
	api.Register(whitelist.NewRPCWhitelist())
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
   --rpc_whitelist 	Comma separated whitelist of allowed services for RPC calls [$MICRO_RPC_WHITELIST]
```

In this case the usage would be

```
micro api --rpc_whitelist go.micro.srv.greeter
```
