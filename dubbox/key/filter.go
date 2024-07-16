package key

import "gopkg.inshopline.com/commons/traffic-plugin-dubbo/traffic_dubbo_consts"

const (
	DubboxProviderTraceFilterKey = "dubbox-provider-trace"
	DubboxConsumerTraceFilterKey = "dubbox-consumer-trace"

	DubboxProviderTagFilterKey = "dubbox-provider-tag"
	DubboxConsumerTagFilterKey = "dubbox-consumer-tag"

	DubboxProviderTrafficFilterKey = traffic_dubbo_consts.TrafficProviderFilter
	DubboxConsumerTrafficFilterKey = traffic_dubbo_consts.TrafficConsumerFilter
)
