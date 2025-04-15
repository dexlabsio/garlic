package middleware

type Config struct {
	Cors *CorsConfig `json:"cors" mapstructure:"cors" yaml:"cors"`
}

type CorsConfig struct {
	AllowedHosts   []string `json:"allowed_hosts" mapstructure:"allowed_hosts" yaml:"allowed_hosts"`
	AllowedMethods []string `json:"allowed_methods" mapstructure:"allowed_methods" yaml:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers" mapstructure:"allowed_headers" yaml:"allowed_headers"`
	ExposedHeaders []string `json:"exposed_headers" mapstructure:"exposed_headers" yaml:"exposed_headers"`
}

func CorsConfigDefaults() *CorsConfig {
	return &CorsConfig{
		AllowedHosts:   []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "Content-Length", "Accept-Encoding",
			"X-CSRF-Token", "Authorization", "X-API-KEY",
		},
		ExposedHeaders: []string{"X-Request-Id", "X-Session-Id"},
	}
}

func Defaults() *Config {
	return &Config{
		Cors: CorsConfigDefaults(),
	}
}
