package tag

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/dubbox/common"
	"dubbo.apache.org/dubbo-go/v3/dubbox/key"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"github.com/dubbogo/gost/log/logger"
)

var (
	_ filter.Filter = (*ConsumerTag)(nil)
)

type ConsumerTag struct{}

func (f *ConsumerTag) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	// transfer keys
	keys := transferKeys()
	if len(keys) > 0 {
		for _, k := range keys {
			val := ctx.Value(k)
			if val == nil {
				continue
			}

			if s, ok := val.(string); ok {
				invocation.SetAttachment(k, s)
			} else {
				logger.Warnf("context value with key '%s' must be string type", k)
				continue
			}
		}
	}

	// java tag transfer
	if v := ctx.Value(common.FlagGrayKey); v != nil {
		if s, ok := v.(string); ok && s == javaTagGrayValue {
			invocation.SetAttachment(javaTagKey, javaTagGrayValue)
		}
	}
	if v := ctx.Value(common.FlagStressKey); v != nil {
		if s, ok := v.(string); ok && s == javaTagStressValue {
			invocation.SetAttachment(javaTagKey, javaTagStressValue)
		}
	}

	return invoker.Invoke(ctx, invocation)
}

func (f *ConsumerTag) OnResponse(ctx context.Context, result protocol.Result, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	// do nothing
	return result
}

func newConsumerTag() filter.Filter {
	return &ConsumerTag{}
}

func init() {
	extension.SetFilter(key.DubboxConsumerTagFilterKey, newConsumerTag)
}
