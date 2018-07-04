# CORS Plugin

The CORS plugin enables the configuration of CORS headers when running micro.

## Usage

Register the plugin before building Micro

```
package main

import (
    "github.com/micro/micro/plugin"
    "github.com/saturn-sonic/go-plugins/micro/cors"
)

func init() {
    plugin.Register(cors.NewPlugin())
}
```

## Configuration

### Environment variables

```
CORS_ALLOWED_HEADERS="X-Custom-Header"
CORS_ALLOWED_ORIGINS="*"
CORS_ALLOWED_METHODS="POST"
```

### Command line
```
$ micro api \
    --cors-allowed-headers=X-Custom-Header \
    --cors-allowed-origins=someotherdomain.com \
    --cors-allowed-methods=POST
```
