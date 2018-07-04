# Disable RPC Plugin

This plugin returns a 403 for /rpc. Nothing more.

## Usage

Register the plugin before building Micro

```
package main

import (
	"github.com/micro/micro/plugin"
	rpc "github.com/saturn-sonic/go-plugins/micro/disable_rpc"
)

func init() {
	plugin.Register(rpc.NewPlugin())
}
```

### Scoped to API

If you like to only apply the plugin for a specific component you can register it with that specifically. 
For example, below you'll see the plugin registered with the API.

```
package main

import (
	"github.com/micro/micro/api"
	rpc "github.com/saturn-sonic/go-plugins/micro/disable_rpc"
)

func init() {
	api.Register(rpc.NewPlugin())
}
```
