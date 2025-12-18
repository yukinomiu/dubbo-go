package trace

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/dubbox/key"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"encoding/json"
	"fmt"
	"github.com/dubbogo/gost/log/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"gopkg.inshopline.com/commons/constx"
	"gopkg.inshopline.com/commons/logx"
	"gopkg.inshopline.com/commons/tracex"
)

var (
	_ filter.Filter = (*ConsumerTrace)(nil)
)

type ConsumerTrace struct{}

func (f *ConsumerTrace) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	serviceKey := invoker.GetURL().ServiceKey()
	providerApplicationName := invoker.GetURL().GetParam("application", serviceKey)
	newCtx, span := clientStartSpan(ctx, serviceKey, providerApplicationName, invocation.MethodName(), invocation.Attachments())
	if newCtx != nil {
		ctx = newCtx
	}
	if span != nil {
		defer func() {
			span.End()
		}()
	}

	// request event
	func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Errorf("client trace filter (request event) panic recovered: %+v", e)
			}
		}()

		var requestEntity string
		if jsonBytes, err := json.Marshal(invocation.Arguments()); err == nil {
			requestEntity = string(jsonBytes)
		} else {
			requestEntity = fmt.Sprintf("%+v", invocation.Arguments())
			logger.Warnf("client side request json serialization error: %s, entity: %s", err.Error(), requestEntity)
		}

		span.AddEvent(EventNameRequest, trace.WithAttributes(attribute.String("entity", requestEntity)))
	}()

	// call
	result := invoker.Invoke(ctx, invocation)
	err := result.Error()

	// response event
	func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Errorf("client trace filter (response event) panic recovered: %+v", e)
			}
		}()

		var responseEntity string
		if jsonBytes, e := json.Marshal(result.Result()); e == nil {
			responseEntity = string(jsonBytes)
		} else {
			responseEntity = fmt.Sprintf("%+v", result.Result())
			logger.Warnf("cliet side response json serialization error: %s, entity: %s", e.Error(), responseEntity)
		}

		span.AddEvent(EventNameResponse, trace.WithAttributes(attribute.String("entity", responseEntity)))
	}()

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

func clientStartSpan(ctx context.Context, serviceKey string, providerApplicationName string, method string, attachment map[string]any) (context.Context, trace.Span) {
	defer func() {
		if e := recover(); e != nil {
			logger.Errorf("client trace start span panic recovered: %+v", e)
		}
	}()

	var (
		ok                = false
		pkgName           string
		simpleServiceName string
		appName           = providerApplicationName
	)

	if ok, pkgName, simpleServiceName, _, _ = getInfoFromServiceKey(serviceKey); !ok {
		// fallback
		pkgName = serviceKey
		simpleServiceName = serviceKey
	}

	tracer := otel.Tracer(constx.COMPONENT_DUBBOX_CONSUMER)
	ctx, span := tracer.Start(
		ctx,
		simpleServiceName+"/"+method,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			semconv.RPCSystemApacheDubbo,
			semconv.RPCServiceKey.String(appName),
			semconv.RPCMethodKey.String(simpleServiceName+"/"+method),
			attribute.String("rpc.package", pkgName),
		),
	)
	ctx = setConsumerSpanAttrs(ctx, span)

	// parent application name
	member, err := baggage.NewMember(constx.SPAN_ATTR_ENV_PARENT_APPLICATION_NAME, constx.ApplicationName())
	if err != nil {
		logger.Warn(ctx, "trace baggage new member error", logx.WithKey("error", err.Error()))
	} else {
		bag := baggage.FromContext(ctx)
		bag, err := bag.SetMember(member)
		if err != nil {
			logger.Warn(ctx, "trace baggage set member error", logx.WithKey("error", err.Error()))
		} else {
			ctx = baggage.ContextWithBaggage(ctx, bag)
		}
	}

	carrier := newAttachmentCarrier(attachment)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return ctx, span
}

func setConsumerSpanAttrs(ctx context.Context, span trace.Span) context.Context {
	// gray
	if flag := ctx.Value(constx.SPAN_ATTR_CICD_VERSION); flag != nil {
		span.SetAttributes(attribute.Key(constx.SPAN_ATTR_CICD_VERSION).String(flag.(string)))
	} else {
		span.SetAttributes(attribute.Key(constx.SPAN_ATTR_CICD_VERSION).String(""))
	}

	// other
	tracex.InjectToSpan(ctx, span)

	return ctx
}

func init() {
	extension.SetFilter(key.DubboxConsumerTraceFilterKey, newConsumerTrace)
}
