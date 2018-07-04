# IP Whitelist Plugin

The IP whitelist plugin is a straight forward plugin for micro which whitelists IP addresses that can access the API.

Current implementation accepts individual IPs or a CIDR.

## Usage

Register the plugin before building Micro

```
package main

import (
	"github.com/micro/micro/plugin"
	ip "github.com/saturn-sonic/go-plugins/micro/ip_whitelist"
)

func init() {
	plugin.Register(ip.NewIPWhitelist())
}
```

It can then be applied on the command line like so.

```
micro --ip_whitelist=10.1.1.10,10.1.1.11,10.1.2.0/24 api
```

### Scoped to API

If you like to only apply the plugin for a specific component you can register it with that specifically. 
For example, below you'll see the plugin registered with the API.

```
package main

import (
	"github.com/micro/micro/api"
	ip "github.com/saturn-sonic/go-plugins/micro/ip_whitelist"
)

func init() {
	api.Register(ip.NewIPWhitelist())
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
   --ip_whitelist 	Comma separated whitelist of allowed IP addresses [$MICRO_IP_WHITELIST]
```

In this case the usage would be

```
micro api --ip_whitelist 10.0.0.0/8
```
