package crypto

type Config struct {
	MasterKey string `json:"masterkey" mapstructure:"masterkey" yaml:"masterkey"`
}

func Defaults() *Config {
	return &Config{
		MasterKey: "INSECURE_CHANGEME",
	}
}
