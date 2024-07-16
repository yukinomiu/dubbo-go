package tag

import "dubbo.apache.org/dubbo-go/v3/dubbox/common"

const (
	javaTagKey         = "dubbo.tag"
	javaTagGrayValue   = "gray"
	javaTagStressValue = "flpt"
)

func transferKeys() []string {
	return []string{
		common.FlagGrayKey,
		common.FlagStressKey,
	}
}
