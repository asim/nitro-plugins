module github.com/micro/go-plugins/registry/etcd/v3

go 1.15

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/micro/go-micro/v3 v3.0.0-beta.2
	github.com/mitchellh/hashstructure v1.0.0
	go.uber.org/zap v1.16.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
