# stats_auth Plugin

stats_auth plugin enables basic auth on the /stats endpoint  

## Usage

Register the plugin before building Micro  

```
package main

import (
	"github.com/micro/micro/plugin"
	"github.com/micro/go-plugins/api/stats_auth"
)

func init() {
	plugin.Register(stats_auth.New())
}
```

You can then set the appropriate variables through command line like so:  

```
micro --enable_stats --stats_auth_user=root --stats_auth_pass=admin --stats_auth_realm=A\ realm\ of\ fun\ and\ happiness api
```

### Scoped to API

If you like to only apply the plugin for a specific component you can register it with that specifically.
For example, below you'll see the plugin registered with the API.  

```
package main

import (
	"github.com/micro/micro/api"
	"github.com/micro/go-plugins/api/stats_auth"
)

func init() {
	api.Register(stats_auth.New())
}
```

Here's the help output:

```
	 --stats_auth_user 								Username used for basic auth for /stats endpoint [$STATS_AUTH_USER]
   --stats_auth_pass 								Password used for basic auth for /stats endpoint [$STATS_AUTH_PASS]
   --stats_auth_realm 							Realm used for basic auth for /stats endpoint. Escape spaces to add multiple words. Optional. Defaults to Access to stats is restricted [$STATS_AUTH_REALM]
```

In this case the usage would be

```
micro --enable_stats api --stats_auth_user=root --stats_auth_pass=admin --stats_auth_realm=A\ realm\ of\ fun\ and\ happiness
```
