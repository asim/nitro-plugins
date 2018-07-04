# Trace/UUID Plugin

The trace/uuid plugin adds a X-Micro-Trace-Id header with a uuid if it does not exist.

## Usage

Register the plugin before building Micro

```
package main

import (
	"github.com/micro/micro/plugin"
	"github.com/saturn-sonic/go-plugins/micro/trace/uuid"
)

func init() {
	plugin.Register(uuid.New())
}
```

### Scoped to API

If you like to only apply the plugin for a specific component you can register it with that specifically. 
For example, below you'll see the plugin registered with the API.

```
package main

import (
	"github.com/micro/micro/api"
	"github.com/saturn-sonic/go-plugins/micro/trace/uuid"
)

func init() {
	api.Register(uuid.New())
}
```
