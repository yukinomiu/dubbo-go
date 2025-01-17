package config

type DubboxConfig struct {
	Flag          Flag          `yaml:"flag" json:"flag" property:"flag"`
	SessionParams SessionParams `yaml:"session-params" json:"session-params" property:"session-params"`
}

type (
	Flag struct {
		DisableProviderCloseConnAfterPingFailed bool `yaml:"disable-provider-close-conn-after-ping-failed" json:"disable-provider-close-conn-after-ping-failed" property:"disable-provider-close-conn-after-ping-failed"`
		DisableConsumerCloseConnAfterPingFailed bool `yaml:"disable-consumer-close-conn-after-ping-failed" json:"disable-consumer-close-conn-after-ping-failed" property:"disable-consumer-close-conn-after-ping-failed"`
	}

	SessionParams struct {
		MaxMsgLen int `yaml:"max-msg-len" json:"max-msg-len" property:"max-msg-len"`
	}
)
