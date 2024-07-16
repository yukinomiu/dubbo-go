package meta

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/dubbox/key"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
)

var (
	_ filter.Filter = (*ConsumerMeta)(nil)
)

const (
	AttachmentKeyConsumerAppName = "remote.application"
)

type ConsumerMeta struct{}

func (f *ConsumerMeta) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	consumerAppName := appName()
	if consumerAppName != "" {
		invocation.SetAttachment(AttachmentKeyConsumerAppName, consumerAppName)
	}

	return invoker.Invoke(ctx, invocation)
}

func (f *ConsumerMeta) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	// do nothing
	return result
}

func appName() string {
	return config.GetApplicationConfig().Name
}

func newConsumerMeta() filter.Filter {
	return &ConsumerMeta{}
}

func init() {
	extension.SetFilter(key.DubboxConsumerMetaFilterKey, newConsumerMeta)
}
