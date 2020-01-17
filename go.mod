module github.com/micro/go-plugins

go 1.13

replace github.com/micro/go-plugins/micro/cors => ./micro/cors

replace github.com/micro/go-plugins/micro/index => ./micro/index

replace github.com/micro/go-plugins/micro/stats_auth => ./micro/stats_auth

replace github.com/micro/go-plugins/micro/ip_whitelist => ./micro/ip_whitelist

replace github.com/micro/go-plugins/micro/gzip => ./micro/gzip

replace github.com/micro/go-plugins/micro/metadata => ./micro/metadata

replace github.com/micro/go-plugins/micro/router => ./micro/router

replace github.com/micro/go-plugins/micro/metrics => ./micro/metrics

replace github.com/micro/go-plugins/micro/metrics/prometheus => ./micro/metrics/prometheus

replace github.com/micro/go-plugins/micro/whitelist => ./micro/whitelist

replace github.com/micro/go-plugins/micro/disable_rpc => ./micro/disable_rpc

replace github.com/micro/go-plugins/micro/auth => ./micro/auth

replace github.com/micro/go-plugins/micro/trace/awsxray => ./micro/trace/awsxray

replace github.com/micro/go-plugins/micro/trace/uuid => ./micro/trace/uuid

replace github.com/micro/go-plugins/micro/header => ./micro/header

replace github.com/micro/go-plugins/client/selector/static => ./client/selector/static

replace github.com/micro/go-plugins/client/selector/shard => ./client/selector/shard

replace github.com/micro/go-plugins/client/selector/label => ./client/selector/label

replace github.com/micro/go-plugins/client/http => ./client/http

replace github.com/micro/go-plugins/codec/jsonrpc2 => ./codec/jsonrpc2

replace github.com/micro/go-plugins/codec/msgpackrpc => ./codec/msgpackrpc

replace github.com/micro/go-plugins/codec/bsonrpc => ./codec/bsonrpc

replace github.com/micro/go-plugins/wrapper/endpoint => ./wrapper/endpoint

replace github.com/micro/go-plugins/wrapper/select/roundrobin => ./wrapper/select/roundrobin

replace github.com/micro/go-plugins/wrapper/select/shard => ./wrapper/select/shard

replace github.com/micro/go-plugins/wrapper/select/version => ./wrapper/select/version

replace github.com/micro/go-plugins/wrapper/breaker/hystrix => ./wrapper/breaker/hystrix

replace github.com/micro/go-plugins/wrapper/breaker/gobreaker => ./wrapper/breaker/gobreaker

replace github.com/micro/go-plugins/wrapper/service => ./wrapper/service

replace github.com/micro/go-plugins/wrapper/ratelimiter/uber => ./wrapper/ratelimiter/uber

replace github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit => ./wrapper/ratelimiter/ratelimit

replace github.com/micro/go-plugins/wrapper/validator => ./wrapper/validator

replace github.com/micro/go-plugins/wrapper/monitoring/victoriametrics => ./wrapper/monitoring/victoriametrics

replace github.com/micro/go-plugins/wrapper/monitoring/prometheus => ./wrapper/monitoring/prometheus

replace github.com/micro/go-plugins/wrapper/trace/awsxray => ./wrapper/trace/awsxray

replace github.com/micro/go-plugins/wrapper/trace/datadog => ./wrapper/trace/datadog

replace github.com/micro/go-plugins/wrapper/trace/opencensus => ./wrapper/trace/opencensus

replace github.com/micro/go-plugins/wrapper/trace/opentracing => ./wrapper/trace/opentracing

replace github.com/micro/go-plugins/service/kubernetes => ./service/kubernetes

replace github.com/micro/go-plugins/web/kubernetes => ./web/kubernetes

replace github.com/micro/go-plugins/server/http => ./server/http

replace github.com/micro/go-plugins/transport/utp => ./transport/utp

replace github.com/micro/go-plugins/transport/grpc => ./transport/grpc

replace github.com/micro/go-plugins/transport/memory => ./transport/memory

replace github.com/micro/go-plugins/transport/quic => ./transport/quic

replace github.com/micro/go-plugins/transport/tcp => ./transport/tcp

replace github.com/micro/go-plugins/transport/rabbitmq => ./transport/rabbitmq

replace github.com/micro/go-plugins/transport/nats => ./transport/nats

replace github.com/micro/go-plugins/config/source/grpc => ./config/source/grpc

replace github.com/micro/go-plugins/config/source/runtimevar => ./config/source/runtimevar

replace github.com/micro/go-plugins/config/source/consul => ./config/source/consul

replace github.com/micro/go-plugins/config/source/mucp => ./config/source/mucp

replace github.com/micro/go-plugins/config/source/vault => ./config/source/vault

replace github.com/micro/go-plugins/config/source/pkger => ./config/source/pkger

replace github.com/micro/go-plugins/config/source/configmap => ./config/source/configmap

replace github.com/micro/go-plugins/config/source/url => ./config/source/url

replace github.com/micro/go-plugins/config/encoder/cue => ./config/encoder/cue

replace github.com/micro/go-plugins/broker/mqtt => ./broker/mqtt

replace github.com/micro/go-plugins/broker/stan => ./broker/stan

replace github.com/micro/go-plugins/broker/googlepubsub => ./broker/googlepubsub

replace github.com/micro/go-plugins/broker/http => ./broker/http

replace github.com/micro/go-plugins/broker/redis => ./broker/redis

replace github.com/micro/go-plugins/broker/sqs => ./broker/sqs

replace github.com/micro/go-plugins/broker/grpc => ./broker/grpc

replace github.com/micro/go-plugins/broker/memory => ./broker/memory

replace github.com/micro/go-plugins/broker/stomp => ./broker/stomp

replace github.com/micro/go-plugins/broker/service => ./broker/service

replace github.com/micro/go-plugins/broker/snssqs => ./broker/snssqs

replace github.com/micro/go-plugins/broker/rabbitmq => ./broker/rabbitmq

replace github.com/micro/go-plugins/broker/kafka => ./broker/kafka

replace github.com/micro/go-plugins/broker/gocloud => ./broker/gocloud

replace github.com/micro/go-plugins/broker/nsq => ./broker/nsq

replace github.com/micro/go-plugins/broker/proxy => ./broker/proxy

replace github.com/micro/go-plugins/broker/nats => ./broker/nats

replace github.com/micro/go-plugins/sync/leader/consul => ./sync/leader/consul

replace github.com/micro/go-plugins/sync/lock/redis => ./sync/lock/redis

replace github.com/micro/go-plugins/sync/lock/consul => ./sync/lock/consul

replace github.com/micro/go-plugins/store/redis => ./store/redis

replace github.com/micro/go-plugins/store/memcached => ./store/memcached

replace github.com/micro/go-plugins/store/consul => ./store/consul

replace github.com/micro/go-plugins/registry/etcdv3 => ./registry/etcdv3

replace github.com/micro/go-plugins/registry/kubernetes => ./registry/kubernetes

replace github.com/micro/go-plugins/registry/cache => ./registry/cache

replace github.com/micro/go-plugins/registry/multi => ./registry/multi

replace github.com/micro/go-plugins/registry/gossip => ./registry/gossip

replace github.com/micro/go-plugins/registry/memory => ./registry/memory

replace github.com/micro/go-plugins/registry/mdns => ./registry/mdns

replace github.com/micro/go-plugins/registry/etcd => ./registry/etcd

replace github.com/micro/go-plugins/registry/service => ./registry/service

replace github.com/micro/go-plugins/registry/consul => ./registry/consul

replace github.com/micro/go-plugins/registry/eureka => ./registry/eureka

replace github.com/micro/go-plugins/registry/proxy => ./registry/proxy

replace github.com/micro/go-plugins/registry/nats => ./registry/nats

replace github.com/micro/go-plugins/registry/zookeeper => ./registry/zookeeper

replace github.com/micro/go-plugins/agent/command/geocode => ./agent/command/geocode

replace github.com/micro/go-plugins/agent/command/whereareyou => ./agent/command/whereareyou

replace github.com/micro/go-plugins/agent/command/animate => ./agent/command/animate

replace github.com/micro/go-plugins/proxy/http => ./proxy/http
