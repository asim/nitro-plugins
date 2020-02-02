module github.com/micro/go-plugins/micro/metrics/v2

go 1.13

require (
	github.com/micro/cli/v2 v2.1.1
	github.com/micro/go-micro/v2 v2.0.0
	github.com/micro/go-plugins/micro/metrics/prometheus/v2 v2.0.0-20200201134528-a2e9624d8e31
	github.com/micro/micro/v2 v2.0.0
)

replace github.com/micro/go-plugins => ../../

replace github.com/micro/go-plugins/v2 => ../../
