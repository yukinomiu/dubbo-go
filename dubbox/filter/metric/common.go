package metric

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
