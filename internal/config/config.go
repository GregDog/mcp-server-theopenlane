package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	DefaultBaseURL  = "https://api.theopenlane.io"
	DefaultLogLevel = "info"
)

// Config is loaded from the process environment.
type Config struct {
	APIToken       string
	BaseURL        string
	OrganizationID string
	LogLevel       string
}

// FromEnv loads configuration. OPENLANE_API_TOKEN is required.
func FromEnv() (Config, error) {
	cfg := Config{
		APIToken:       strings.TrimSpace(os.Getenv("OPENLANE_API_TOKEN")),
		BaseURL:        strings.TrimSpace(os.Getenv("OPENLANE_BASE_URL")),
		OrganizationID: strings.TrimSpace(os.Getenv("OPENLANE_ORGANIZATION_ID")),
		LogLevel:       strings.TrimSpace(os.Getenv("OPENLANE_MCP_LOG_LEVEL")),
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = DefaultLogLevel
	}
	if cfg.APIToken == "" {
		return Config{}, fmt.Errorf("OPENLANE_API_TOKEN is required")
	}
	return cfg, nil
}
