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
	_ filter.Filter = (*ProviderTag)(nil)
)

type ProviderTag struct{}

func (f *ProviderTag) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	attach := invocation.Attachments()
	if len(attach) > 0 {
		// transfer keys
		keys := transferKeys()
		if len(keys) > 0 {
			for _, k := range keys {
				if val, ok := attach[k]; ok {
					if s, ok := val.(string); ok {
						ctx = context.WithValue(ctx, k, s)
					} else {
						logger.Debugf("attachment value with key '%s' must be string type", k)
					}
				}
			}
		}

		// java tag transfer
		if v, ok := attach[javaTagKey]; ok {
			if s, ok := v.(string); ok {
				switch s {
				case javaTagGrayValue: // gray
					ctx = context.WithValue(ctx, common.FlagGrayKey, javaTagGrayValue)
					break
				case javaTagStressValue: // stress
					ctx = context.WithValue(ctx, common.FlagStressKey, javaTagStressValue)
					break
				}
			} else {
				logger.Debugf("attachment value with key '%s' must be string type", javaTagKey)
			}
		}
	}

	return invoker.Invoke(ctx, invocation)
}

func (f *ProviderTag) OnResponse(ctx context.Context, result protocol.Result, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	// do nothing
	return result
}

func newProviderTag() filter.Filter {
	return &ProviderTag{}
}

func init() {
	extension.SetFilter(key.DubboxProviderTagFilterKey, newProviderTag)
}
