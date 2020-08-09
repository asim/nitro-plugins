# IP Allow Plugin

The IP allow plugin is a straight forward plugin for micro which allows IP addresses that can allow the API.

Current implementation accepts individual IPs or a CIDR.

## Usage

Register the plugin before building Micro

```
package main

import (
	"github.com/micro/micro/plugin"
	ip "github.com/micro/go-plugins/api/ip_allow"
)

func init() {
	plugin.Register(ip.NewIPAllow())
}
```

It can then be applied on the command line like so.

```
micro --ip_allow=10.1.1.10,10.1.1.11,10.1.2.0/24 api
```

### Scoped to API

If you like to only apply the plugin for a specific component you can register it with that specifically. 
For example, below you'll see the plugin registered with the API.

```
package main

import (
	"github.com/micro/micro/api"
	ip "github.com/micro/go-plugins/api/ip_allow"
)

func init() {
	api.Register(ip.NewIPAllow())
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
   --ip_allow 	Comma separated list of allowed IP addresses [$MICRO_IP_ALLOW]
```

In this case the usage would be

```
micro api --ip_allow 10.0.0.0/8
```
