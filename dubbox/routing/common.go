package routing

import "dubbo.apache.org/dubbo-go/v3/dubbox/common"

func tagRoutingKeys() []string {
	return []string{
		common.FlagGrayKey,
		common.FlagStressKey,
	}
}
