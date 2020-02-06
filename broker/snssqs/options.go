package snssqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/micro/go-micro/v2/broker"
)

type maxMessagesKey struct{}
type sqsConfigKey struct{}
type snsConfigKey struct{}
type stsConfigKey struct{}

// MaxReceiveMessages indicates how many messages a receive operation should pull
// during any single call
func MaxReceiveMessages(max int64) broker.SubscribeOption {
	return setSubscribeOption(maxMessagesKey{}, max)
}

type visibilityTimeoutKey struct{}

// VisibilityTimeout controls how long a message is hidden from other queue consumers
// before being put back. If a consumer does not delete the message, it will be put back
// even if it was "processed"
func VisibilityTimeout(seconds int64) broker.SubscribeOption {
	return setSubscribeOption(visibilityTimeoutKey{}, seconds)
}

type waitTimeSecondsKey struct{}

// WaitTimeSeconds controls the length of long polling for available messages
func WaitTimeSeconds(seconds int64) broker.SubscribeOption {
	return setSubscribeOption(waitTimeSecondsKey{}, seconds)
}

type validateOnPublishKey struct{}

// ValidateOnPublish determines whether to pre-validate messages before they're published
// This has a significant performance impact
func ValidateOnPublish(validate bool) broker.PublishOption {
	return setPublishOption(validateOnPublishKey{}, validate)
}

// SNSConfig add AWS config options to the sns client
func SNSConfig(c *aws.Config) broker.Option {
	return setBrokerOption(snsConfigKey{}, c)
}

// SQSConfig add AWS config options to the sqs client
func SQSConfig(c *aws.Config) broker.Option {
	return setBrokerOption(sqsConfigKey{}, c)
}

func STSConfig(c *aws.Config) broker.Option {
	return setBrokerOption(stsConfigKey{}, c)
}
