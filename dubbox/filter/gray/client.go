package metric

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/dubbox/key"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"github.com/dubbogo/gost/log/logger"
	"gopkg.inshopline.com/commons/grayx/allgrayx"
)

var (
	_ filter.Filter = (*ConsumerGray)(nil)
)

type ConsumerGray struct{}

func (f *ConsumerGray) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	if shouldBeGray(ctx) {
		ctx = toGrayCtx(ctx)
	}

	result := invoker.Invoke(ctx, invocation)
	return result
}

func (f *ConsumerGray) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	// do nothing
	return result
}

func newConsumerGray() filter.Filter {
	return &ConsumerGray{}
}

func shouldBeGray(ctx context.Context) bool {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("get full gray flag from gray-x panic recovered: %+v", r)
		}
	}()

	return allgrayx.IsAllGray(ctx)
}

func toGrayCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, allgrayx.GRAY_KEY, allgrayx.GRAY_VALUE)
}

func init() {
	extension.SetFilter(key.DubboxConsumerGrayFilterKey, newConsumerGray)
}
