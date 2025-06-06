package config

// AppConfig holds application configuration
type AppConfig struct {
	ServerPort string
	ServerHost string
}

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() AppConfig {
	return AppConfig{
		ServerPort: "8080",
		ServerHost: "0.0.0.0",
	}
}
