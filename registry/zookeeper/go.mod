module github.com/micro/go-plugins/registry/zookeeper/v3

go 1.15

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/go-zookeeper/zk v1.0.2
	github.com/micro/go-micro/v3 v3.0.0-beta.2
	github.com/mitchellh/hashstructure v1.0.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
)
