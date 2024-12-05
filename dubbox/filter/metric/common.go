package metric

import (
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"fmt"
)

var (
	defaultClientMetricMap = map[string]string{"side": "consumer", "uri_tag": "c", "component_type": "dubbox"}
	defaultServerMetricMap = map[string]string{"side": "provider", "uri_tag": "s", "component_type": "dubbox"}
)

func errToCode(err error) string {
	if err == nil {
		return "success"
	}

	return err.Error()
}

func getInterface(invoker protocol.Invoker) string {
	return invoker.GetURL().Interface()
}

func metricKey(dubboInterface, method string) string {
	return fmt.Sprintf("dubbo/%s/%s", dubboInterface, method)
}
