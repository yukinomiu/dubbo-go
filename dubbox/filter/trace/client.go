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
	_ filter.Filter = (*ConsumerTrace)(nil)
)

type ConsumerTrace struct{}

func (f *ConsumerTrace) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	newCtx, span := clientStartSpan(ctx, invoker.GetURL().ServiceKey(), invocation.MethodName(), invocation.Attachments())
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

func (f *ConsumerTrace) OnResponse(_ context.Context, result protocol.Result, _ protocol.Invoker, _ protocol.Invocation) protocol.Result {
	// do nothing
	return result
}

func newConsumerTrace() filter.Filter {
	return &ConsumerTrace{}
}

func clientStartSpan(ctx context.Context, serviceName, method string, attachment map[string]any) (context.Context, trace.Span) {
	defer func() {
		if e := recover(); e != nil {
			logger.Errorf("client trace start span panic recovered: %v", e)
		}
	}()

	tracer := otel.Tracer(constx.COMPONENT_DUBBOX_CONSUMER)
	ctx, span := tracer.Start(
		ctx,
		method,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			semconv.RPCSystemApacheDubbo,
			semconv.RPCServiceKey.String(serviceName),
			semconv.RPCMethodKey.String(method),
			SpanMetricKey.String(method),
		),
	)
	ctx = setConsumerSpanAttrs(ctx, method, span)

	carrier := newAttachmentCarrier(attachment)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return ctx, span
}

func setConsumerSpanAttrs(ctx context.Context, method string, span trace.Span) context.Context {
	// application name
	span.SetAttributes(attribute.Key(constx.SPAN_ATTR_ENV_APPLICATION_NAME).String(constx.ApplicationName()))

	// parent application name
	member, err := baggage.NewMember(constx.SPAN_ATTR_ENV_PARENT_APPLICATION_NAME, constx.ApplicationName())
	if err != nil {
		logger.Warnf("trace baggage new member error: %v", err)
	} else {
		bag := baggage.FromContext(ctx)
		bag, err := bag.SetMember(member)
		if err != nil {
			logger.Warnf("trace baggage set member error: %v", err)
		} else {
			ctx = baggage.ContextWithBaggage(ctx, bag)
		}
	}

	// middle ware detail
	span.SetAttributes(attribute.Key(constx.SPAN_ATTR_MIDDLEWARE_TYPE).String(SpanAttrMiddlewareType))
	span.SetAttributes(attribute.Key(constx.SPAN_ATTR_MIDDLEWARE_ACTION).String(SpanAttrMiddlewareConsumer))

	// gray
	if flag := ctx.Value(constx.SPAN_ATTR_CICD_VERSION); flag != nil {
		span.SetAttributes(attribute.Key(constx.SPAN_ATTR_CICD_VERSION).String(flag.(string)))
	} else {
		span.SetAttributes(attribute.Key(constx.SPAN_ATTR_CICD_VERSION).String(""))
	}

	// metric key
	span.SetAttributes(attribute.Key(constx.SPAN_ATTR_METRICX_KEY).String(method))

	return ctx
}

func init() {
	extension.SetFilter(key.DubboxConsumerTraceFilterKey, newConsumerTrace)
}
