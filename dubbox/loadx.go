package dubbox

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	"gopkg.inshopline.com/commons/loadx"
)

func RegisterOrderedLoadX(starter loadx.Starter) {
	loadx.RegisterComponent(loadx.NewSimpleComponent(&loadx.Descriptor{
		ComponentInfo: loadx.ComponentInfo{
			Name: "dubbox",
			Desc: "dubbox based on apache dubbo-go",
		},
		Ordered: true,
		Order:   450,
		Starter: starter,
		Stopper: func() error {
			config.BeforeShutdown()
			return nil
		},
	}))
}
