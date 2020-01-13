module github.com/micro/go-plugins

go 1.13

replace k8s.io/api => k8s.io/api v0.0.0-20190708174958-539a33f6e817

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d

replace k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190708180123-608cd7da68f7

replace k8s.io/client-go => k8s.io/client-go v11.0.0+incompatible

replace k8s.io/component-base => k8s.io/component-base v0.0.0-20190708175518-244289f83105

replace github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.0

require (
	cloud.google.com/go/pubsub v1.0.1
	cuelang.org/go v0.0.14
	github.com/DataDog/datadog-go v3.2.0+incompatible // indirect
	github.com/Shopify/sarama v1.24.1
	github.com/VictoriaMetrics/metrics v1.9.2
	github.com/abbot/go-http-auth v0.4.1-0.20181019201920-860ed7f246ff
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/anacrolix/sync v0.2.0 // indirect
	github.com/anacrolix/utp v0.0.0-20180219060659-9e0e1d1d0572
	github.com/asim/go-awsxray v0.0.0-20161209120537-0d8a60b6e205
	github.com/asim/go-bson v0.0.0-20160318195205-84522947cabd
	github.com/aws/aws-sdk-go v1.25.31
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/clbanning/x2j v0.0.0-20191024224557-825249438eec // indirect
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/franela/goblin v0.0.0-20181003173013-ead4ad1d2727 // indirect
	github.com/franela/goreq v0.0.0-20171204163338-bcd34c9993f8 // indirect
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-log/log v0.1.0
	github.com/go-redsync/redsync v1.3.1
	github.com/go-stomp/stomp v2.0.3+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gorilla/websocket v1.4.1
	github.com/hashicorp/consul/api v1.2.0
	github.com/hashicorp/memberlist v0.1.5
	github.com/hashicorp/vault/api v1.0.4
	github.com/hudl/fargo v1.3.0
	github.com/juju/ratelimit v1.0.1
	github.com/markbates/pkger v0.12.4
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.18.1-0.20191216151820-1ea6390eaedd
	github.com/micro/micro v1.18.1-0.20191216152715-bc3ebd75c9ff
	github.com/minio/highwayhash v1.0.0
	github.com/mitchellh/hashstructure v1.0.0
	github.com/nats-io/nats-streaming-server v0.16.2 // indirect
	github.com/nats-io/nats.go v1.9.1
	github.com/nats-io/stan.go v0.5.0
	github.com/nsqio/go-nsq v1.0.7
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.2.1
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4
	github.com/rs/cors v1.7.0
	github.com/samuel/go-zookeeper v0.0.0-20190923202752-2cc03de413da
	github.com/sirupsen/logrus v1.4.2
	github.com/sony/gobreaker v0.4.1
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	github.com/stretchr/testify v1.4.0
	github.com/tinylib/msgp v1.1.0
	go.opencensus.io v0.22.2
	go.uber.org/ratelimit v0.1.0
	gocloud.dev v0.17.0
	gocloud.dev/pubsub/rabbitpubsub v0.17.0
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6
	golang.org/x/text v0.3.2
	google.golang.org/api v0.13.0
	google.golang.org/genproto v0.0.0-20191108220845-16a3f7862a1a
	google.golang.org/grpc v1.25.1
	gopkg.in/DataDog/dd-trace-go.v1 v1.19.0
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20160220154919-db14e161995a // indirect
	gopkg.in/gcfg.v1 v1.2.3 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ldap.v3 v3.1.0
	gopkg.in/redis.v3 v3.6.4
	k8s.io/api v0.0.0-20191109101513-0171b7c15da1
	k8s.io/apimachinery v0.0.0-20191111054156-6eb29fdf75dc
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20191030222137-2b95a09bc58d // indirect
)
