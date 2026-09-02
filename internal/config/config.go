package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultBaseURL       = "https://api.theopenlane.io"
	DefaultLogLevel      = "info"
	DefaultTransport     = "stdio"
	DefaultHTTPAddr      = ":8080"
	DefaultHTTPMaxBody   = 32 * 1024 * 1024 // 32 MiB
	DefaultMaxUpload     = 10 * 1024 * 1024 // 10 MiB
	DefaultUploadTimeout = 2 * time.Minute
)

// Config is loaded from the process environment.
type Config struct {
	APIToken         string
	BaseURL          string
	OrganizationID   string
	LogLevel         string
	AllowWrite       bool
	AllowDelete      bool
	Transport        string
	HTTPAddr         string
	HTTPJSON         bool
	HTTPMaxBodyBytes int64
	MaxUploadBytes   int64
	UploadTimeout    time.Duration
}

// FromEnv loads configuration. OPENLANE_API_TOKEN is required.
func FromEnv() (Config, error) {
	cfg := Config{
		APIToken:         strings.TrimSpace(os.Getenv("OPENLANE_API_TOKEN")),
		BaseURL:          strings.TrimSpace(os.Getenv("OPENLANE_BASE_URL")),
		OrganizationID:   strings.TrimSpace(os.Getenv("OPENLANE_ORGANIZATION_ID")),
		LogLevel:         strings.TrimSpace(os.Getenv("OPENLANE_MCP_LOG_LEVEL")),
		AllowWrite:       parseBoolEnv(os.Getenv("OPENLANE_ALLOW_WRITE")),
		AllowDelete:      parseBoolEnv(os.Getenv("OPENLANE_ALLOW_DELETE")),
		Transport:        strings.TrimSpace(os.Getenv("OPENLANE_MCP_TRANSPORT")),
		HTTPAddr:         strings.TrimSpace(os.Getenv("OPENLANE_MCP_HTTP_ADDR")),
		HTTPJSON:         parseBoolEnv(os.Getenv("OPENLANE_MCP_HTTP_JSON")),
		HTTPMaxBodyBytes: parseInt64Env(os.Getenv("OPENLANE_MCP_HTTP_MAX_BODY_BYTES"), DefaultHTTPMaxBody),
		MaxUploadBytes:   parseInt64Env(os.Getenv("OPENLANE_MCP_MAX_UPLOAD_BYTES"), DefaultMaxUpload),
		UploadTimeout:    parseDurationEnv(os.Getenv("OPENLANE_MCP_UPLOAD_TIMEOUT"), DefaultUploadTimeout),
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = DefaultLogLevel
	}
	if cfg.Transport == "" {
		cfg.Transport = DefaultTransport
	}
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = DefaultHTTPAddr
	}
	if cfg.APIToken == "" {
		return Config{}, fmt.Errorf("OPENLANE_API_TOKEN is required")
	}
	switch strings.ToLower(cfg.Transport) {
	case "stdio", "http":
	default:
		return Config{}, fmt.Errorf("OPENLANE_MCP_TRANSPORT must be stdio or http")
	}
	return cfg, nil
}

func parseInt64Env(value string, fallback int64) int64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

func parseDurationEnv(value string, fallback time.Duration) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	d, err := time.ParseDuration(value)
	if err != nil || d <= 0 {
		return fallback
	}
	return d
}

func parseBoolEnv(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
