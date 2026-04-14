package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	AuditLogPath string            `mapstructure:"audit_log_path"`
	Services     map[string]Service `mapstructure:"services"`
}

// Service represents a configured external service.
type Service struct {
	Name    string `mapstructure:"name"`
	KeyFile string `mapstructure:"key_file"`
	RotateDays int `mapstructure:"rotate_days"`
}

// Load reads the configuration from disk using viper.
func Load(cfgFile string) (*Config, error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("resolving home dir: %w", err)
		}
		viper.AddConfigPath(filepath.Join(home, ".fencepost"))
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("FENCEPOST")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}

	if cfg.AuditLogPath == "" {
		cfg.AuditLogPath = defaultAuditLogPath()
	}

	return &cfg, nil
}

func defaultAuditLogPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".fencepost", "audit.log")
}
