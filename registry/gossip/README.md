# Gossip Registry

Gossip is a registry plugin for go-micro which uses hashicorp/memberlist to broadcast registry information 
via the SWIM protocol.

## Usage

Import the plugin as per usual

```go
import _ "github.com/saturn-sonic/go-plugins/registry/gossip"
```

Start with the registry flag

```go
go run service.go --registry=gossip
```

On startup you'll see something like

```go
2016/06/19 14:05:43 Local memberlist node 127.0.0.1:45465
```

To join this gossip ring use `--registry=gossip --registry_address 127.0.0.1:45465` when starting other nodes
