module github.com/micro/go-plugins/service/kubernetes/v2

go 1.13

require (
	github.com/micro/go-micro/v2 v2.0.0
	github.com/micro/go-plugins/client/selector/static/v2 v2.0.0-00010101000000-000000000000
	github.com/micro/go-plugins/registry/kubernetes/v2 v2.0.0-00010101000000-000000000000
)

replace github.com/micro/go-plugins/client/selector/static/v2 => ../../client/selector/static

replace github.com/micro/go-plugins/registry/kubernetes/v2 => ../../registry/kubernetes
