package traffic

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/dubbox/hook"
	"github.com/dubbogo/gost/log/logger"
	_ "gopkg.inshopline.com/commons/traffic-plugin-dubbo/traffic_dubbo" // register traffic filter
	trafficFilter "gopkg.inshopline.com/commons/traffic-plugin-dubbo/traffic_dubbo"
	"strconv"
)

func init() {
	hook.RegisterAfterRCLoadHook(func() {
		rc := config.GetRootConfig()
		if rc == nil {
			logger.Warn("dubbox: can not find root config")
			return
		}
		if !isProvider(rc) {
			logger.Info("dubbox: no provider found, skip set traffic filter port")
			return
		}

		var portStr *string
		if len(rc.Protocols) > 0 {
			for _, pc := range rc.Protocols {
				if pc != nil && pc.Name == "dubbo" {
					portStr = &pc.Port
					break
				}
			}
		}

		if portStr == nil {
			logger.Warn("dubbox: can not find dubbo protocol port")
			return
		}

		port, err := strconv.Atoi(*portStr)
		if err != nil {
			logger.Errorf("dubbox: convert port string to int failed: %v, port string: %s",
				err, *portStr)
			return
		}
		trafficFilter.SetProviderPort(int64(port))
		logger.Infof("dubbox: init traffic record plugin with port: %d", port)
	})
}

func isProvider(rc *config.RootConfig) bool {
	if rc != nil &&
		rc.Provider != nil &&
		len(rc.Provider.Services) > 0 {
		return true
	}

	return false
}
