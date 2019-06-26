module github.com/micro/go-plugins

go 1.12

require (
	cloud.google.com/go v0.40.0
	github.com/Shopify/sarama v1.22.1
	github.com/abbot/go-http-auth v0.4.1-0.20181019201920-860ed7f246ff
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/anacrolix/envpprof v1.0.0 // indirect
	github.com/anacrolix/missinggo v1.1.0 // indirect
	github.com/anacrolix/sync v0.0.0-20180808010631-44578de4e778 // indirect
	github.com/anacrolix/utp v0.0.0-20180219060659-9e0e1d1d0572
	github.com/asim/go-awsxray v0.0.0-20161209120537-0d8a60b6e205
	github.com/asim/go-bson v0.0.0-20160318195205-84522947cabd
	github.com/aws/aws-sdk-go v1.19.44
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20190329173943-551aad21a668
	github.com/bradfitz/iter v0.0.0-20190303215204-33e6a9893b0c // indirect
	github.com/cenkalti/backoff v2.1.1+incompatible // indirect
	github.com/clbanning/x2j v0.0.0-20180326210544-5e605d46809c // indirect
	github.com/coreos/bbolt v1.3.2 // indirect
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/franela/goblin v0.0.0-20181003173013-ead4ad1d2727 // indirect
	github.com/franela/goreq v0.0.0-20171204163338-bcd34c9993f8 // indirect
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/go-log/log v0.1.0
	github.com/go-redsync/redsync v1.2.0
	github.com/go-stomp/stomp v2.0.3+incompatible
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/golang/protobuf v1.3.1
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/gorilla/websocket v1.4.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/hashicorp/vault/api v1.0.2
	github.com/hudl/fargo v1.2.0
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/juju/ratelimit v1.0.1
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.7.0
	github.com/micro/micro v1.5.0
	github.com/minio/highwayhash v1.0.0
	github.com/mitchellh/hashstructure v1.0.0
	github.com/nats-io/nats-streaming-server v0.15.1 // indirect
	github.com/nats-io/nats.go v1.8.1
	github.com/nats-io/stan.go v0.5.0
	github.com/nsqio/go-nsq v1.0.7
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90
	github.com/rs/cors v1.6.0
	github.com/samuel/go-zookeeper v0.0.0-20180130194729-c4fab1ac1bec
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/sony/gobreaker v0.4.1
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94
	github.com/stretchr/testify v1.3.0
	github.com/tinylib/msgp v1.1.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.opencensus.io v0.22.0
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/ratelimit v0.0.0-20180316092928-c15da0234277
	go.uber.org/zap v1.10.0 // indirect
	gocloud.dev v0.15.0
	gocloud.dev/pubsub/rabbitpubsub v0.15.0
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/text v0.3.2
	google.golang.org/api v0.6.0
	google.golang.org/genproto v0.0.0-20190620144150-6af8c5fc6601
	google.golang.org/grpc v1.21.1
	gopkg.in/DataDog/dd-trace-go.v1 v1.14.0
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20160220154919-db14e161995a // indirect
	gopkg.in/gcfg.v1 v1.2.3 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ldap.v3 v3.0.3
	gopkg.in/redis.v3 v3.6.4
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/utils v0.0.0-20190529001817-6999998975a7 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)
