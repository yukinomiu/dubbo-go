package key

import "gopkg.inshopline.com/commons/traffic-plugin-dubbo/traffic_dubbo_consts"

const (
	DubboxProviderTraceFilterKey = "dubbox-provider-trace"
	DubboxConsumerTraceFilterKey = "dubbox-consumer-trace"

	DubboxProviderMetricFilterKey = "dubbox-provider-metric"
	DubboxConsumerMetricFilterKey = "dubbox-consumer-metric"

	DubboxProviderTagFilterKey = "dubbox-provider-tag"
	DubboxConsumerTagFilterKey = "dubbox-consumer-tag"

	DubboxConsumerMetaFilterKey = "dubbox-consumer-meta"

	DubboxProviderSentinelFilterKey = "dubbox-provider-sentinel"
	DubboxCircuitBreakerFilterKey   = "dubbox-circuit-breaker"

	DubboxProviderPropagationFilterKey = "dubbox-provider-propagation"

	DubboxConsumerChaosPluginFilterKey = "dubbox-consumer-chaos-plugin"

	DubboxConsumerGrayFilterKey = "dubbox-consumer-gray"

	DubboxProviderTrafficFilterKey = traffic_dubbo_consts.TrafficProviderFilter
	DubboxConsumerTrafficFilterKey = traffic_dubbo_consts.TrafficConsumerFilter
)
