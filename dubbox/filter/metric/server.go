package metric

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/dubbox/key"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"github.com/dubbogo/gost/log/logger"
	"gopkg.inshopline.com/commons/metricx"
	"time"
)

var (
	_ filter.Filter = (*ProviderMetric)(nil)
)

type ProviderMetric struct{}

func (f *ProviderMetric) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	start := time.Now()
	result := invoker.Invoke(ctx, invocation)
	duration := time.Since(start)

	var err error
	if result != nil {
		err = result.Error()
	}
	serverRecordMetric(invocation.MethodName(), duration, err)
	return result
}

func (f *ProviderMetric) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	// do nothing
	return result
}

func newProviderMetric() filter.Filter {
	return &ProviderMetric{}
}

func serverRecordMetric(methodName string, duration time.Duration, err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Errorf("server metric panic recovered: %v", e)
		}
	}()

	durationMs := duration.Milliseconds()
	metricx.CountDurationAndCode(methodName, durationMs, errToCode(err), defaultServerMetricMap)
	return
}

func init() {
	extension.SetFilter(key.DubboxProviderMetricFilterKey, newProviderMetric)
}
