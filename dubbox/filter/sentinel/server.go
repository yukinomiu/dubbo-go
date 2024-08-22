package sentinel

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/dubbox/key"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/dubbogo/gost/log/logger"
	"gopkg.inshopline.com/commons/sentinel-go/flow"
)

var (
	_ filter.Filter = (*ProviderSentinel)(nil)
)

type ProviderSentinel struct{}

func (f *ProviderSentinel) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	var (
		interfaceName = invocation.GetAttachmentWithDefaultValue(constant.InterfaceKey, "")
		methodName    = invocation.MethodName()
		resourceName  = interfaceName + ":" + methodName
	)

	entry, blockErr := flow.Entry(
		resourceName,
		sentinel.WithTrafficType(base.Inbound),
		sentinel.WithResourceType(base.ResTypeRPC),
	)
	if blockErr != nil {
		logger.Warnf("dubbo call was blocked by sentinel, resource: %s, method: %s",
			resourceName, methodName)

		result := &protocol.RPCResult{}
		result.SetError(blockErr)
		return result
	}

	defer entry.Exit()
	return invoker.Invoke(ctx, invocation)
}

func (f *ProviderSentinel) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	// do nothing
	return result
}

func newProviderSentinel() filter.Filter {
	return &ProviderSentinel{}
}

func init() {
	extension.SetFilter(key.DubboxProviderSentinelFilterKey, newProviderSentinel)
}
