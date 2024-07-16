package chaosplugin

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/dubbox/key"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"encoding/json"
	"fmt"
	"github.com/dubbogo/gost/log/logger"
	"gopkg.inshopline.com/commons/chaos-go-agent/chaosgo"
)

var (
	_ filter.Filter = (*ConsumerChaosPlugin)(nil)
)

type ConsumerChaosPlugin struct{}

func (f *ConsumerChaosPlugin) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	// bypass generic invocation
	if invocation.IsGenericInvocation() {
		return invoker.Invoke(ctx, invocation)
	}

	// chaos plugin inject
	interfaceName := invoker.GetURL().GetParam(constant.InterfaceKey, "")
	method := invocation.MethodName()

	replyReplaced := chaosInject(ctx, interfaceName, method, invocation.Reply())
	if replyReplaced {
		result := &protocol.RPCResult{}
		result.SetResult(invocation.Reply())
		return result
	}

	return invoker.Invoke(ctx, invocation)
}

func (f *ConsumerChaosPlugin) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	// do nothing
	return result
}

func chaosInject(ctx context.Context, serverServiceName, method string, reply interface{}) (replyReplaced bool) {
	if !chaosgo.CheckInjectable(ctx) {
		logger.Debugf("client chaos inject disabled. method: %s#%s", serverServiceName, method)
		return
	}

	// inject
	filePath := serverServiceName + "#" + method
	value := chaosgo.TryInject(chaosgo.GetFailPointName(clientModuleName, filePath))
	if value != nil {
		if s, ok := value.(string); ok {
			if err := json.Unmarshal([]byte(s), reply); err == nil {
				replyReplaced = true
				logger.Debugf("client chaos inject reply replaced, file path: %s, JSON: %s", filePath, s)
			} else {
				// replyReplaced = false
				logger.Errorf("client chaos inject JSON unmarshal error, file path: %s, JSON: %s, err: %v", filePath, s, err)
			}
		}
	}

	return
}

func newConsumerChaosPlugin() filter.Filter {
	// initiate chaos plugin
	if ret := chaosgo.Start(); ret != "success" {
		panic(fmt.Sprintf("initiate chaos plugin failed: %s", ret))
	}

	return &ConsumerChaosPlugin{}
}

func init() {
	extension.SetFilter(key.DubboxConsumerChaosPluginFilterKey, newConsumerChaosPlugin)
}
