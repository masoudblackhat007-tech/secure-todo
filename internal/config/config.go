package config

type AppConfig struct {
	HTTPPort string
}

func Load() AppConfig {
	// TODO: later - read from env
	return AppConfig{
		HTTPPort: "8080",
	}
}
