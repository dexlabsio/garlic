package httpclient

type Config struct {
	URL string `mapstructure:"url" yaml:"url"`
}

func Defaults() *Config {
	return &Config{
		URL: "http://localhost",
	}
}
