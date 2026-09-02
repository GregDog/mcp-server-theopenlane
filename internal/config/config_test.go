package config

import (
	"testing"
)

func TestFromEnvRequiresToken(t *testing.T) {
	t.Setenv("OPENLANE_API_TOKEN", "")
	t.Setenv("OPENLANE_BASE_URL", "")
	t.Setenv("OPENLANE_ORGANIZATION_ID", "")
	t.Setenv("OPENLANE_MCP_LOG_LEVEL", "")
	t.Setenv("OPENLANE_ALLOW_WRITE", "")
	t.Setenv("OPENLANE_ALLOW_DELETE", "")

	if _, err := FromEnv(); err == nil {
		t.Fatal("expected error when token is missing")
	}
}

func TestFromEnvDefaults(t *testing.T) {
	t.Setenv("OPENLANE_API_TOKEN", "tola_example")
	t.Setenv("OPENLANE_BASE_URL", "")
	t.Setenv("OPENLANE_ORGANIZATION_ID", "")
	t.Setenv("OPENLANE_MCP_LOG_LEVEL", "")
	t.Setenv("OPENLANE_ALLOW_WRITE", "")
	t.Setenv("OPENLANE_ALLOW_DELETE", "")

	cfg, err := FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != DefaultBaseURL {
		t.Fatalf("base url: got %q", cfg.BaseURL)
	}
	if cfg.LogLevel != DefaultLogLevel {
		t.Fatalf("log level: got %q", cfg.LogLevel)
	}
	if cfg.APIToken != "tola_example" {
		t.Fatalf("token not loaded")
	}
	if cfg.AllowWrite || cfg.AllowDelete {
		t.Fatal("expected write/delete disabled by default")
	}
}

func TestFromEnvOverrides(t *testing.T) {
	t.Setenv("OPENLANE_API_TOKEN", "tolp_example")
	t.Setenv("OPENLANE_BASE_URL", "https://openlane.example.internal")
	t.Setenv("OPENLANE_ORGANIZATION_ID", "org_123")
	t.Setenv("OPENLANE_MCP_LOG_LEVEL", "debug")
	t.Setenv("OPENLANE_ALLOW_WRITE", "")
	t.Setenv("OPENLANE_ALLOW_DELETE", "")

	cfg, err := FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != "https://openlane.example.internal" {
		t.Fatalf("base url: got %q", cfg.BaseURL)
	}
	if cfg.OrganizationID != "org_123" {
		t.Fatalf("org: got %q", cfg.OrganizationID)
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("log level: got %q", cfg.LogLevel)
	}
}

func TestFromEnvAllowWrite(t *testing.T) {
	t.Setenv("OPENLANE_API_TOKEN", "tola_example")
	t.Setenv("OPENLANE_BASE_URL", "")
	t.Setenv("OPENLANE_ORGANIZATION_ID", "")
	t.Setenv("OPENLANE_MCP_LOG_LEVEL", "")
	t.Setenv("OPENLANE_ALLOW_WRITE", "true")
	t.Setenv("OPENLANE_ALLOW_DELETE", "")

	cfg, err := FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.AllowWrite {
		t.Fatal("expected allow write true")
	}
	if cfg.AllowDelete {
		t.Fatal("expected allow delete false")
	}
}

func TestFromEnvAllowDelete(t *testing.T) {
	t.Setenv("OPENLANE_API_TOKEN", "tola_example")
	t.Setenv("OPENLANE_BASE_URL", "")
	t.Setenv("OPENLANE_ORGANIZATION_ID", "")
	t.Setenv("OPENLANE_MCP_LOG_LEVEL", "")
	t.Setenv("OPENLANE_ALLOW_WRITE", "")
	t.Setenv("OPENLANE_ALLOW_DELETE", "true")

	cfg, err := FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AllowWrite {
		t.Fatal("expected allow write false")
	}
	if !cfg.AllowDelete {
		t.Fatal("expected allow delete true")
	}
}
