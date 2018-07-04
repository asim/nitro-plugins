# Header Plugin

The header plugin allows you to set http headers returned to the client

## Usage

Create a plugin file

```
package main

import (
	"github.com/saturn-sonic/go-plugins/micro/header"
	"github.com/micro/micro/plugin"
)

func init() {
	plugin.Register(header.NewPlugin())
}
```

Build micro binary with the plugin

It can then be flagged as so

```
micro --header "Access-Control-Allow-Headers=User-Agent,X-Requested-With,If-Modified-Since" api
```
