package config

type DubboxConfig struct {
	Flag Flag `yaml:"flag" json:"flag,omitempty" property:"flag"`
}

type (
	Flag struct {
		DisableProviderCloseConnAfterPingFailed bool `yaml:"disable-provider-close-conn-after-ping-failed" json:"disable-provider-close-conn-after-ping-failed" property:"disable-provider-close-conn-after-ping-failed"`
		DisableConsumerCloseConnAfterPingFailed bool `yaml:"disable-consumer-close-conn-after-ping-failed" json:"disable-consumer-close-conn-after-ping-failed" property:"disable-consumer-close-conn-after-ping-failed"`
	}
)
