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
	"gopkg.inshopline.com/commons/tracex"
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

	// request event
	func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Errorf("server trace filter (request event) panic recovered: %+v", e)
			}
		}()

		var requestEntity string
		if jsonBytes, err := json.Marshal(invocation.Arguments()); err == nil {
			requestEntity = string(jsonBytes)
		} else {
			requestEntity = fmt.Sprintf("%+v", invocation.Arguments())
			logger.Warnf("server side request json serialization error: %s, entity: %s", err.Error(), requestEntity)
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
				logger.Errorf("server trace filter (response event) panic recovered: %+v", e)
			}
		}()

		var responseEntity string
		if jsonBytes, e := json.Marshal(result.Result()); e == nil {
			responseEntity = string(jsonBytes)
		} else {
			responseEntity = fmt.Sprintf("%+v", result.Result())
			logger.Warnf("server side response json serialization error: %s, entity: %s", e.Error(), responseEntity)
		}

		span.AddEvent(EventNameResponse, trace.WithAttributes(attribute.String("entity", responseEntity)))
	}()

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

func serverStartSpan(ctx context.Context, serviceKey string, method string, attachment map[string]any) (context.Context, trace.Span) {
	defer func() {
		if e := recover(); e != nil {
			logger.Errorf("server trace start span panic recovered: %+v", e)
		}
	}()

	var (
		ok                = false
		pkgName           string
		simpleServiceName string
		appName           = constx.ApplicationName()
	)

	if ok, pkgName, simpleServiceName, _, _ = getInfoFromServiceKey(serviceKey); !ok {
		// fallback
		pkgName = serviceKey
		simpleServiceName = serviceKey
	}

	carrier := newAttachmentCarrier(attachment)
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	bag := baggage.FromContext(ctx)
	ctx = baggage.ContextWithBaggage(ctx, bag)
	spanCtx := trace.SpanContextFromContext(ctx)

	tracer := otel.Tracer(constx.COMPONENT_DUBBOX_PROVIDER)
	ctx, span := tracer.Start(
		trace.ContextWithRemoteSpanContext(ctx, spanCtx),
		simpleServiceName+"/"+method,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			semconv.RPCSystemApacheDubbo,
			semconv.RPCServiceKey.String(appName),
			semconv.RPCMethodKey.String(simpleServiceName+"/"+method),
			attribute.String("rpc.package", pkgName),
		),
	)
	ctx = setProviderSpanAttrs(ctx, attachment, span)

	return ctx, span
}

func setProviderSpanAttrs(ctx context.Context, attachment map[string]any, span trace.Span) context.Context {
	// gray
	if flag, ok := attachment[constx.SPAN_ATTR_CICD_VERSION]; ok && flag != nil {
		span.SetAttributes(attribute.Key(constx.SPAN_ATTR_CICD_VERSION).String(flag.(string)))
	} else {
		span.SetAttributes(attribute.Key(constx.SPAN_ATTR_CICD_VERSION).String(""))
	}

	// parent application name
	bag := baggage.FromContext(ctx)
	span.SetAttributes(attribute.Key(constx.SPAN_ATTR_ENV_PARENT_APPLICATION_NAME).String(bag.Member(constx.SPAN_ATTR_ENV_PARENT_APPLICATION_NAME).Value()))

	// other
	tracex.InjectToSpan(ctx, span)

	return ctx
}

func init() {
	extension.SetFilter(key.DubboxProviderTraceFilterKey, newProviderTrace)
}
