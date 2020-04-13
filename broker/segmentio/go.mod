module github.com/micro/go-plugins/broker/segmentio/v2

go 1.13

require (
	github.com/frankban/quicktest v1.4.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v2 v2.3.0
	github.com/micro/go-plugins/broker/kafka/v2 v2.3.0
	github.com/micro/go-plugins/codec/segmentio/v2 v2.3.0
	github.com/pierrec/lz4 v2.2.6+incompatible // indirect
	github.com/segmentio/kafka-go v0.3.5
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
)

replace github.com/micro/go-plugins/codec/segmentio/v2 => ../../codec/segmentio
