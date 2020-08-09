# GZIP Plugin

The gzip plugin is a plugin for the micro toolkit which enables gzipping of http response

## Usage

Register the plugin before building Micro

```
package main

import (
	"github.com/micro/micro/plugin"
	"github.com/micro/go-plugins/api/gzip"
)

func init() {
	plugin.Register(gzip.New())
}
```

### Scoped to API

If you like to only apply the plugin for a specific component you can register it with that specifically. 
For example, below you'll see the plugin registered with the API.

```
package main

import (
	"github.com/micro/micro/api"
	"github.com/micro/go-plugins/api/gzip"
)

func init() {
	api.Register(gzip.New())
}
```
