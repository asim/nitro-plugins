module github.com/micro/go-plugins

require (
	cloud.google.com/go v0.44.0
	github.com/Azure/azure-pipeline-go v0.2.2 // indirect
	github.com/Azure/azure-storage-blob-go v0.7.0 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.6.0 // indirect
	github.com/DataDog/zstd v1.4.1 // indirect
	github.com/Shopify/sarama v1.23.1
	github.com/abbot/go-http-auth v0.4.1-0.20181019201920-860ed7f246ff
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/anacrolix/envpprof v1.0.0 // indirect
	github.com/anacrolix/missinggo v1.1.1 // indirect
	github.com/anacrolix/sync v0.0.0-20180808010631-44578de4e778 // indirect
	github.com/anacrolix/utp v0.0.0-20180219060659-9e0e1d1d0572
	github.com/asim/go-awsxray v0.0.0-20161209120537-0d8a60b6e205
	github.com/asim/go-bson v0.0.0-20160318195205-84522947cabd
	github.com/aws/aws-sdk-go v1.23.0
	github.com/bradfitz/gomemcache v0.0.0-20190329173943-551aad21a668
	github.com/bradfitz/iter v0.0.0-20190303215204-33e6a9893b0c // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/clbanning/x2j v0.0.0-20180326210544-5e605d46809c // indirect
	github.com/coreos/etcd v3.3.17+incompatible
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/franela/goblin v0.0.0-20181003173013-ead4ad1d2727 // indirect
	github.com/franela/goreq v0.0.0-20171204163338-bcd34c9993f8 // indirect
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/go-log/log v0.1.0
	github.com/go-redsync/redsync v1.3.0
	github.com/go-stomp/stomp v2.0.3+incompatible
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.3.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190430165422-3e4dfb77656c // indirect
	github.com/gorilla/websocket v1.4.1
	github.com/grpc-ecosystem/grpc-gateway v1.9.5 // indirect
	github.com/hashicorp/consul/api v1.2.0
	github.com/hashicorp/go-hclog v0.9.2 // indirect
	github.com/hashicorp/memberlist v0.1.5
	github.com/hashicorp/raft v1.1.1 // indirect
	github.com/hashicorp/vault v1.1.4
	github.com/huandu/xstrings v1.2.0 // indirect
	github.com/hudl/fargo v1.2.0
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/juju/ratelimit v1.0.1
	github.com/lib/pq v1.2.0 // indirect
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.12.0
	github.com/micro/micro v1.11.0
	github.com/minio/highwayhash v1.0.0
	github.com/mitchellh/hashstructure v1.0.0
	github.com/nats-io/nats-streaming-server v0.15.1 // indirect
	github.com/nats-io/nats.go v1.8.1
	github.com/nats-io/stan.go v0.5.0
	github.com/nsqio/go-nsq v1.0.7
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90
	github.com/rcrowley/go-metrics v0.0.0-20190706150252-9beb055b7962 // indirect
	github.com/rs/cors v1.6.0
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20180130194729-c4fab1ac1bec
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/assertions v1.0.1 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190710185942-9d28bd7c0945 // indirect
	github.com/sony/gobreaker v0.4.1
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94
	github.com/stretchr/testify v1.4.0
	github.com/tinylib/msgp v1.1.0
	go.opencensus.io v0.22.0
	go.uber.org/ratelimit v0.1.0
	gocloud.dev v0.17.0
	gocloud.dev/pubsub/rabbitpubsub v0.15.0
	golang.org/x/net v0.0.0-20191011234655-491137f69257
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/text v0.3.2
	google.golang.org/api v0.8.0
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
	google.golang.org/grpc v1.24.0
	gopkg.in/DataDog/dd-trace-go.v1 v1.16.1
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20160220154919-db14e161995a // indirect
	gopkg.in/gcfg.v1 v1.2.3 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
	gopkg.in/jcmturner/gokrb5.v7 v7.3.0 // indirect
	gopkg.in/ldap.v3 v3.0.3
	gopkg.in/redis.v3 v3.6.4
	k8s.io/api v0.0.0-20190726022912-69e1bce1dad5
	k8s.io/apimachinery v0.0.0-20190727130956-f97a4e5b4abc
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/klog v0.3.3 // indirect
	k8s.io/utils v0.0.0-20190712204705-3dccf664f023 // indirect
)

exclude (
	github.com/Sirupsen/logrus v1.1.0
	github.com/Sirupsen/logrus v1.1.1
	github.com/Sirupsen/logrus v1.2.0
	github.com/Sirupsen/logrus v1.3.0
	github.com/Sirupsen/logrus v1.4.0
	github.com/Sirupsen/logrus v1.4.1
	github.com/Sirupsen/logrus v1.4.2
	github.com/gogo/protobuf v0.0.0-00010101000000-000000000000
	sourcegraph.com/sourcegraph/go-diff v0.5.1
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.2.0+incompatible
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.1
	github.com/golang/lint => github.com/golang/lint v0.0.0-20190227174305-8f45f776aaf1
	github.com/hashicorp/vault/api => github.com/hashicorp/vault/sdk v1.0.2
	github.com/hashicorp/vault/sdk => github.com/hashicorp/vault/sdk v0.1.8
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.4
	github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.4.0
	github.com/ugorji/go/codec => github.com/ugorji/go/codec v0.0.0-20181204163529-d75b2dcb6bc8
	k8s.io/api => k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go => k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)

go 1.13
