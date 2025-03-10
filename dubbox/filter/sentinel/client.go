package sentinel

import (
	"context"
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
	_ filter.Filter = (*ClientCircuitBreaker)(nil)
)

type ClientCircuitBreaker struct{}

func (f *ClientCircuitBreaker) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	var (
		interfaceName = invoker.GetURL().Service()
		methodName    = invocation.MethodName()
		resourceName  = circuitBreakerResourceNamePrefix + interfaceName + ":" + methodName
	)
	entry, blockErr := flow.Entry(
		resourceName,
		sentinel.WithTrafficType(base.Outbound),
		sentinel.WithResourceType(base.ResTypeRPC),
	)
	if blockErr != nil {
		logger.Warnf("dubbo call was blocked by circuit breaker, resource: %s, method: %s, block: %s",
			resourceName, methodName, blockErr.Error())
		result := &protocol.RPCResult{}
		result.SetError(blockErr)
		return result
	}

	defer entry.Exit()
	result := invoker.Invoke(ctx, invocation)
	if result != nil && result.Error() != nil {
		sentinel.TraceError(entry, result.Error())
	}
	return result
}

func (f *ClientCircuitBreaker) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	// do nothing
	return result
}

func newClientCircuitBreaker() filter.Filter {
	return &ClientCircuitBreaker{}
}

func init() {
	extension.SetFilter(key.DubboxCircuitBreakerFilterKey, newClientCircuitBreaker)
}
