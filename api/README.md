# Micro Plugins

This directory containers plugins for the micro toolkit itself. These plugins differ from 
go-micro plugins in that they allow adding functionality to the API, CLI, Sidecar, etc.

Learn more about micro plugins at [github.com/micro/micro/plugin](https://github.com/micro/micro/tree/master/plugin).

## The Plugin Interface

```go
// Plugin is the interface for plugins to micro. It differs from go-micro in that it's for
// the micro API, Web, Sidecar, CLI. It's a method of building middleware for the HTTP side.
type Plugin interface {
    // Global Flags
    Flags() []cli.Flag
    // Sub-commands
    Commands() []cli.Command
    // Handle is the middleware handler for HTTP requests. We pass in
    // the existing handler so it can be wrapped to create a call chain.
    Handler() Handler
    // Init called when command line args are parsed.
    // The initialised cli.Context is passed in.
    Init(*cli.Context) error
    // Name of the plugin
    String() string
}

// Manager is the plugin manager which stores plugins and allows them to be retrieved.
// This is used by all the components of micro.
type Manager interface {
        Plugins() map[string]Plugin
        Register(name string, plugin Plugin) error
}

// Handler is the plugin middleware handler which wraps an existing http.Handler passed in.
// Its the responsibility of the Handler to call the next http.Handler in the chain.
type Handler func(http.Handler) http.Handler
```
