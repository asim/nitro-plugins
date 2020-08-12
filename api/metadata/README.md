# Metadata Plugin

The metadata plugin lets you set metadata for any micro command or to inject headers into the api, sidecar and web proxy

## Usage

Create a plugin file

```
package main

import (
	"github.com/micro/go-plugins/api/metadata"
	"github.com/micro/micro/plugin"
)

func init() {
	plugin.Register(metadata.NewPlugin())
}
```

Build micro binary with the plugin

It can then be flagged as so

```
micro --metadata foo=bar --metadata bar=baz query go.micro.srv.greeter Say.Hello '{"name": "john"}'
```
