package routing

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"github.com/dubbogo/gost/log/logger"
)

func SetRoutingTag(ctx context.Context, invocation protocol.Invocation) {
	keys := tagRoutingKeys()
	if len(keys) > 0 {
		for _, k := range keys {
			val := ctx.Value(k)
			if val == nil {
				continue
			}

			if s, ok := val.(string); ok {
				invocation.SetAttachment(constant.Tagkey, s)
				logger.Debugf("set tag-routing tag: %s", s)
				break
			} else {
				logger.Warnf("context value with tag-routing key '%s' must be string type", k)
				continue
			}
		}
	}
}
