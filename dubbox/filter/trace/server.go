package trace

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/dubbox/key"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"github.com/dubbogo/gost/log/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"gopkg.inshopline.com/commons/constx"
)

var (
	_ filter.Filter = (*ProviderTrace)(nil)
)

type ProviderTrace struct{}

func (f *ProviderTrace) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	newCtx, span := serverStartSpan(ctx, invoker.GetURL().ServiceKey(), invocation.MethodName(), invocation.Attachments())
	if newCtx != nil {
		ctx = newCtx
	}
	if span != nil {
		defer func() {
			span.End()
		}()
	}

	result := invoker.Invoke(ctx, invocation)
	err := result.Error()
	setSpanStatus(span, err)
	return result
}

func (f *ProviderTrace) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	// do nothing
	return result
}

func newProviderTrace() filter.Filter {
	return &ProviderTrace{}
}

func serverStartSpan(ctx context.Context, serviceName, method string, attachment map[string]any) (context.Context, trace.Span) {
	defer func() {
		if e := recover(); e != nil {
			logger.Errorf("server trace start span panic recovered: %v", e)
		}
	}()

	carrier := newAttachmentCarrier(attachment)
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	bag := baggage.FromContext(ctx)
	ctx = baggage.ContextWithBaggage(ctx, bag)
	spanCtx := trace.SpanContextFromContext(ctx)

	tracer := otel.Tracer(constx.COMPONENT_DUBBOX_PROVIDER)
	ctx, span := tracer.Start(
		trace.ContextWithRemoteSpanContext(ctx, spanCtx),
		method,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			semconv.RPCSystemApacheDubbo,
			semconv.RPCServiceKey.String(serviceName),
			semconv.RPCMethodKey.String(method),
			SpanMetricKey.String(method),
		),
	)
	ctx = setProviderSpanAttrs(ctx, method, attachment, span)

	return ctx, span
}

func setProviderSpanAttrs(ctx context.Context, method string, attachment map[string]any, span trace.Span) context.Context {
	// application name
	span.SetAttributes(attribute.Key(constx.SPAN_ATTR_ENV_APPLICATION_NAME).String(constx.ApplicationName()))

	// parent application name
	bag := baggage.FromContext(ctx)
	span.SetAttributes(attribute.Key(constx.SPAN_ATTR_ENV_PARENT_APPLICATION_NAME).String(bag.Member(constx.SPAN_ATTR_ENV_PARENT_APPLICATION_NAME).Value()))

	// middle ware detail
	span.SetAttributes(attribute.Key(constx.SPAN_ATTR_MIDDLEWARE_TYPE).String(SpanAttrMiddlewareType))
	span.SetAttributes(attribute.Key(constx.SPAN_ATTR_MIDDLEWARE_ACTION).String(SpanAttrMiddlewareProvider))

	// gray
	if flag, ok := attachment[constx.SPAN_ATTR_CICD_VERSION]; ok && flag != nil {
		span.SetAttributes(attribute.Key(constx.SPAN_ATTR_CICD_VERSION).String(flag.(string)))
	} else {
		span.SetAttributes(attribute.Key(constx.SPAN_ATTR_CICD_VERSION).String(""))
	}

	// metric key
	span.SetAttributes(attribute.Key(constx.SPAN_ATTR_METRICX_KEY).String(method))

	return ctx
}

func init() {
	extension.SetFilter(key.DubboxProviderTraceFilterKey, newProviderTrace)
}
