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
	AllowWrite     bool
	AllowDelete    bool
}

// FromEnv loads configuration. OPENLANE_API_TOKEN is required.
func FromEnv() (Config, error) {
	cfg := Config{
		APIToken:       strings.TrimSpace(os.Getenv("OPENLANE_API_TOKEN")),
		BaseURL:        strings.TrimSpace(os.Getenv("OPENLANE_BASE_URL")),
		OrganizationID: strings.TrimSpace(os.Getenv("OPENLANE_ORGANIZATION_ID")),
		LogLevel:       strings.TrimSpace(os.Getenv("OPENLANE_MCP_LOG_LEVEL")),
		AllowWrite:     parseBoolEnv(os.Getenv("OPENLANE_ALLOW_WRITE")),
		AllowDelete:    parseBoolEnv(os.Getenv("OPENLANE_ALLOW_DELETE")),
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

func parseBoolEnv(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
