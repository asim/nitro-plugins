module github.com/micro/go-plugins/config/source/etcd/v3

go 1.15

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/micro/go-micro/v3 v3.0.0-beta.2
)
