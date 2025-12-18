package trace

import (
	"github.com/dubbogo/gost/log/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"gopkg.inshopline.com/commons/tracex"
	"strings"
)

const (
	SpanErrorMessageKey = attribute.Key("error.message")

	EventNameRequest  = "request"
	EventNameResponse = "response"
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
		span.SetAttributes(SpanErrorMessageKey.String(err.Error()))
	}
	tracex.SetSpanStatus(span, err)
}

func getInfoFromServiceKey(serviceKey string) (ok bool, pkgName string, simpleServiceName string, group string, version string) {
	s1 := strings.Split(serviceKey, "/")
	if len(s1) != 2 {
		ok = false
		return
	}

	group = s1[0]
	s2 := strings.Split(s1[1], ":")
	if len(s2) != 2 {
		ok = false
		return
	}

	version = s2[1]
	fullService := s2[0]
	idx := strings.LastIndex(fullService, ".")
	if idx <= 0 || idx >= len(fullService)-1 {
		ok = false
		return
	}

	pkgName = fullService[:idx]
	simpleServiceName = fullService[idx+1:]
	ok = true
	return
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
