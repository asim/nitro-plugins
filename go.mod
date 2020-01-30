module github.com/micro/go-plugins/v2

go 1.13

replace github.com/micro/go-plugins/v2/micro/cors => ./micro/cors

replace github.com/micro/go-plugins/v2/micro/index => ./micro/index

replace github.com/micro/go-plugins/v2/micro/stats_auth => ./micro/stats_auth

replace github.com/micro/go-plugins/v2/micro/ip_whitelist => ./micro/ip_whitelist

replace github.com/micro/go-plugins/v2/micro/gzip => ./micro/gzip

replace github.com/micro/go-plugins/v2/micro/metadata => ./micro/metadata

replace github.com/micro/go-plugins/v2/micro/router => ./micro/router

replace github.com/micro/go-plugins/v2/micro/metrics => ./micro/metrics

replace github.com/micro/go-plugins/v2/micro/metrics/prometheus => ./micro/metrics/prometheus

replace github.com/micro/go-plugins/v2/micro/whitelist => ./micro/whitelist

replace github.com/micro/go-plugins/v2/micro/disable_rpc => ./micro/disable_rpc

replace github.com/micro/go-plugins/v2/micro/auth => ./micro/auth

replace github.com/micro/go-plugins/v2/micro/trace/awsxray => ./micro/trace/awsxray

replace github.com/micro/go-plugins/v2/micro/trace/uuid => ./micro/trace/uuid

replace github.com/micro/go-plugins/v2/micro/header => ./micro/header

replace github.com/micro/go-plugins/v2/client/selector/static => ./client/selector/static

replace github.com/micro/go-plugins/v2/client/selector/shard => ./client/selector/shard

replace github.com/micro/go-plugins/v2/client/selector/label => ./client/selector/label

replace github.com/micro/go-plugins/v2/client/http => ./client/http

replace github.com/micro/go-plugins/v2/codec/jsonrpc2 => ./codec/jsonrpc2

replace github.com/micro/go-plugins/v2/codec/msgpackrpc => ./codec/msgpackrpc

replace github.com/micro/go-plugins/v2/codec/bsonrpc => ./codec/bsonrpc

replace github.com/micro/go-plugins/v2/wrapper/endpoint => ./wrapper/endpoint

replace github.com/micro/go-plugins/v2/wrapper/select/roundrobin => ./wrapper/select/roundrobin

replace github.com/micro/go-plugins/v2/wrapper/select/shard => ./wrapper/select/shard

replace github.com/micro/go-plugins/v2/wrapper/select/version => ./wrapper/select/version

replace github.com/micro/go-plugins/v2/wrapper/breaker/hystrix => ./wrapper/breaker/hystrix

replace github.com/micro/go-plugins/v2/wrapper/breaker/gobreaker => ./wrapper/breaker/gobreaker

replace github.com/micro/go-plugins/v2/wrapper/service => ./wrapper/service

replace github.com/micro/go-plugins/v2/wrapper/ratelimiter/uber => ./wrapper/ratelimiter/uber

replace github.com/micro/go-plugins/v2/wrapper/ratelimiter/ratelimit => ./wrapper/ratelimiter/ratelimit

replace github.com/micro/go-plugins/v2/wrapper/validator => ./wrapper/validator

replace github.com/micro/go-plugins/v2/wrapper/monitoring/victoriametrics => ./wrapper/monitoring/victoriametrics

replace github.com/micro/go-plugins/v2/wrapper/monitoring/prometheus => ./wrapper/monitoring/prometheus

replace github.com/micro/go-plugins/v2/wrapper/trace/awsxray => ./wrapper/trace/awsxray

replace github.com/micro/go-plugins/v2/wrapper/trace/datadog => ./wrapper/trace/datadog

replace github.com/micro/go-plugins/v2/wrapper/trace/opencensus => ./wrapper/trace/opencensus

replace github.com/micro/go-plugins/v2/wrapper/trace/opentracing => ./wrapper/trace/opentracing

replace github.com/micro/go-plugins/v2/service/kubernetes => ./service/kubernetes

replace github.com/micro/go-plugins/v2/web/kubernetes => ./web/kubernetes

replace github.com/micro/go-plugins/v2/server/http => ./server/http

replace github.com/micro/go-plugins/v2/transport/utp => ./transport/utp

replace github.com/micro/go-plugins/v2/transport/grpc => ./transport/grpc

replace github.com/micro/go-plugins/v2/transport/memory => ./transport/memory

replace github.com/micro/go-plugins/v2/transport/quic => ./transport/quic

replace github.com/micro/go-plugins/v2/transport/tcp => ./transport/tcp

replace github.com/micro/go-plugins/v2/transport/rabbitmq => ./transport/rabbitmq

replace github.com/micro/go-plugins/v2/transport/nats => ./transport/nats

replace github.com/micro/go-plugins/v2/config/source/grpc => ./config/source/grpc

replace github.com/micro/go-plugins/v2/config/source/runtimevar => ./config/source/runtimevar

replace github.com/micro/go-plugins/v2/config/source/consul => ./config/source/consul

replace github.com/micro/go-plugins/v2/config/source/mucp => ./config/source/mucp

replace github.com/micro/go-plugins/v2/config/source/vault => ./config/source/vault

replace github.com/micro/go-plugins/v2/config/source/pkger => ./config/source/pkger

replace github.com/micro/go-plugins/v2/config/source/configmap => ./config/source/configmap

replace github.com/micro/go-plugins/v2/config/source/url => ./config/source/url

replace github.com/micro/go-plugins/v2/config/encoder/cue => ./config/encoder/cue

replace github.com/micro/go-plugins/v2/broker/mqtt => ./broker/mqtt

replace github.com/micro/go-plugins/v2/broker/stan => ./broker/stan

replace github.com/micro/go-plugins/v2/broker/googlepubsub => ./broker/googlepubsub

replace github.com/micro/go-plugins/v2/broker/http => ./broker/http

replace github.com/micro/go-plugins/v2/broker/redis => ./broker/redis

replace github.com/micro/go-plugins/v2/broker/sqs => ./broker/sqs

replace github.com/micro/go-plugins/v2/broker/grpc => ./broker/grpc

replace github.com/micro/go-plugins/v2/broker/memory => ./broker/memory

replace github.com/micro/go-plugins/v2/broker/stomp => ./broker/stomp

replace github.com/micro/go-plugins/v2/broker/service => ./broker/service

replace github.com/micro/go-plugins/v2/broker/snssqs => ./broker/snssqs

replace github.com/micro/go-plugins/v2/broker/rabbitmq => ./broker/rabbitmq

replace github.com/micro/go-plugins/v2/broker/kafka => ./broker/kafka

replace github.com/micro/go-plugins/v2/broker/gocloud => ./broker/gocloud

replace github.com/micro/go-plugins/v2/broker/nsq => ./broker/nsq

replace github.com/micro/go-plugins/v2/broker/proxy => ./broker/proxy

replace github.com/micro/go-plugins/v2/broker/nats => ./broker/nats

replace github.com/micro/go-plugins/v2/sync/leader/consul => ./sync/leader/consul

replace github.com/micro/go-plugins/v2/sync/lock/redis => ./sync/lock/redis

replace github.com/micro/go-plugins/v2/sync/lock/consul => ./sync/lock/consul

replace github.com/micro/go-plugins/v2/store/redis => ./store/redis

replace github.com/micro/go-plugins/v2/store/memcached => ./store/memcached

replace github.com/micro/go-plugins/v2/store/consul => ./store/consul

replace github.com/micro/go-plugins/v2/registry/etcdv3 => ./registry/etcdv3

replace github.com/micro/go-plugins/v2/registry/kubernetes => ./registry/kubernetes

replace github.com/micro/go-plugins/v2/registry/cache => ./registry/cache

replace github.com/micro/go-plugins/v2/registry/multi => ./registry/multi

replace github.com/micro/go-plugins/v2/registry/gossip => ./registry/gossip

replace github.com/micro/go-plugins/v2/registry/memory => ./registry/memory

replace github.com/micro/go-plugins/v2/registry/mdns => ./registry/mdns

replace github.com/micro/go-plugins/v2/registry/etcd => ./registry/etcd

replace github.com/micro/go-plugins/v2/registry/service => ./registry/service

replace github.com/micro/go-plugins/v2/registry/consul => ./registry/consul

replace github.com/micro/go-plugins/v2/registry/eureka => ./registry/eureka

replace github.com/micro/go-plugins/v2/registry/proxy => ./registry/proxy

replace github.com/micro/go-plugins/v2/registry/nats => ./registry/nats

replace github.com/micro/go-plugins/v2/registry/zookeeper => ./registry/zookeeper

replace github.com/micro/go-plugins/v2/agent/command/geocode => ./agent/command/geocode

replace github.com/micro/go-plugins/v2/agent/command/whereareyou => ./agent/command/whereareyou

replace github.com/micro/go-plugins/v2/agent/command/animate => ./agent/command/animate

replace github.com/micro/go-plugins/v2/proxy/http => ./proxy/http

require (
	github.com/micro/go-micro/v2 v2.0.0
	github.com/micro/micro/v2 v2.0.0
)
