package trace

import (
	"github.com/dubbogo/gost/log/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	SpanAttrMiddlewareType     = "dubbo"
	SpanAttrMiddlewareConsumer = "consumer"
	SpanAttrMiddlewareProvider = "provider"

	SpanMetricKey    = attribute.Key("metricx-key")
	SpanStatusOkDesc = "OK"
)

var (
	_ propagation.TextMapCarrier = (*attachmentCarrier)(nil)
)

func setSpanStatus(span trace.Span, err error) {
	if span == nil {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			logger.Errorf("set span attribute panic recovered: %v", e)
		}
	}()

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, SpanStatusOkDesc)
	}
}

type attachmentCarrier struct {
	attachment map[string]any
}

func (c *attachmentCarrier) Get(key string) string {
	v, exists := c.attachment[key]
	if v != nil && exists {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

func (c *attachmentCarrier) Set(key string, value string) {
	c.attachment[key] = value
}

func (c *attachmentCarrier) Keys() []string {
	keys := make([]string, 0, len(c.attachment))
	for k, v := range c.attachment {
		if _, ok := v.(string); ok {
			keys = append(keys, k)
		}
	}

	return keys
}

func newAttachmentCarrier(attachment map[string]interface{}) propagation.TextMapCarrier {
	// attachment can not be nil, skip validation in this inner function
	return &attachmentCarrier{
		attachment: attachment,
	}
}
