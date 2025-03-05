package metric

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/dubbox/key"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
)

var (
	_ filter.Filter = (*ProviderPropagation)(nil)
)

type ProviderPropagation struct{}

func (f *ProviderPropagation) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	// dubbox fix: remove provider async key
	delete(invocation.Attachments(), constant.AsyncKey)
	result := invoker.Invoke(ctx, invocation)
	return result
}

func (f *ProviderPropagation) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	// do nothing
	return result
}

func newProviderPropagation() filter.Filter {
	return &ProviderPropagation{}
}

func init() {
	extension.SetFilter(key.DubboxProviderPropagationFilterKey, newProviderPropagation)
}
