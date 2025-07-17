package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// ProviderConfig holds configuration for a specific LLM provider.
type ProviderConfig struct {
	BaseURL string `mapstructure:"BASE_URL"`
	Type    string `mapstructure:"TYPE"`
}

type Config struct {
	DBUser        string `mapstructure:"POSTGRES_USER"`
	DBPassword    string `mapstructure:"POSTGRES_PASSWORD"`
	DBName        string `mapstructure:"POSTGRES_DB"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	ServerPort    string `mapstructure:"SERVER_PORT"`
	EncryptionKey string `mapstructure:"ENCRYPTION_KEY"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	Providers map[string]ProviderConfig `mapstructure:"PROVIDERS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetEnvPrefix("PROVIDER")
	viper.AutomaticEnv()

	requiredEnvs := []string{
		"DB_HOST", "DB_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB",
		"JWT_SECRET", "ENCRYPTION_KEY", "SERVER_PORT",
	}

	for _, env := range requiredEnvs {
		if bindErr := viper.BindEnv(env); bindErr != nil {
			return config, fmt.Errorf("failed to bind env %s: %w", env, bindErr)
		}
		if val := os.Getenv(env); val == "" {
			return config, fmt.Errorf("required environment variable %s is not set", env)
		}
		viper.Set(env, os.Getenv(env))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return
}
