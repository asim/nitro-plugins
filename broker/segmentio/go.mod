module github.com/micro/go-plugins/broker/segmentio/v2

go 1.13

require (
	github.com/Shopify/sarama v1.25.0
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v2 v2.3.0
	github.com/micro/go-plugins/broker/kafka/v2 v2.3.0
	github.com/segmentio/kafka-go v0.3.5
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
)

replace github.com/micro/go-micro/v2 => /home/vtolstov/devel/projects/micro/go-micro
