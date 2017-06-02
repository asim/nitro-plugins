# Cache blacklist Selector

github.com/micro/go-plugins/selector/cache
	+
github.com/micro/go-plugins/selector/blacklist

# Cache
The cache selector is a client side load balancer for go-micro. This selector uses the registry Watcher to cache selected services.
It defaults random hashed strategy for load balancing requests across services.

# Blacklist
The blacklist selector is a go-micro/selector which filters nodes based on which have errored out.
It operates much like a circuit breaker. If a node returns an error 3 consecutive times it will
be blacklisted. After a period of 30 seconds it will be put back into the list of nodes.



```
cli := client.NewClient(
		client.Selector(
			cacheblack.NewSelector(
				selector.Registry(retcdv3),
				cacheblack.CacheTTL(30*time.Second),
				cacheblack.BlackTTL(30),
				cacheblack.BlackErrCount(3),
			),
		),
	)
```
