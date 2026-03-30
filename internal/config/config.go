package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	KeyAccessToken  = "access_token"
	KeyClientID     = "client_id"
	KeyClientSecret = "client_secret"
	KeyAPIKey       = "api_key"
)

func Init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not find home directory: %v\n", err)
		os.Exit(1)
	}

	configDir := filepath.Join(home, ".line-cli")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	_ = viper.ReadInConfig()
}

func Get(key string) string {
	return viper.GetString(key)
}

func Set(key, value string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".line-cli")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	viper.Set(key, value)
	return viper.WriteConfigAs(filepath.Join(configDir, "config.yaml"))
}

func RequireToken() string {
	token := Get(KeyAccessToken)
	if token == "" {
		fmt.Fprintln(os.Stderr, "error: no access token set. Run: line config set --token <token>")
		os.Exit(1)
	}
	return token
}
