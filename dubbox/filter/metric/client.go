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
	_ filter.Filter = (*ConsumerMetric)(nil)
)

type ConsumerMetric struct{}

func (f *ConsumerMetric) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	start := time.Now()
	result := invoker.Invoke(ctx, invocation)
	duration := time.Since(start)

	var err error
	if result != nil {
		err = result.Error()
	}
	clientRecordMetric(invocation.MethodName(), duration, err)
	return result
}

func (f *ConsumerMetric) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	// do nothing
	return result
}

func newConsumerMetric() filter.Filter {
	return &ConsumerMetric{}
}

func clientRecordMetric(methodName string, duration time.Duration, err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Errorf("client metric panic recovered: %v", e)
		}
	}()

	durationMs := duration.Milliseconds()
	metricx.CountDurationAndCode(methodName, durationMs, errToCode(err), defaultClientMetricMap)
	return
}

func init() {
	extension.SetFilter(key.DubboxConsumerMetricFilterKey, newConsumerMetric)
}
