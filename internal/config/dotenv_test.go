package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvFileDoesNotOverrideExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("OPENLANE_API_TOKEN=from_file\nOPENLANE_MCP_LOG_LEVEL=debug\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("OPENLANE_API_TOKEN", "from_env")
	t.Setenv("OPENLANE_BASE_URL", "")
	t.Setenv("OPENLANE_ORGANIZATION_ID", "")
	t.Setenv("OPENLANE_MCP_LOG_LEVEL", "")
	t.Setenv("OPENLANE_ALLOW_WRITE", "")
	t.Setenv("OPENLANE_ALLOW_DELETE", "")
	t.Setenv("OPENLANE_MCP_TRANSPORT", "")
	t.Setenv("OPENLANE_MCP_HTTP_ADDR", "")

	if err := loadDotEnvFile(path); err != nil {
		t.Fatal(err)
	}

	cfg, err := FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIToken != "from_env" {
		t.Fatalf("token: got %q want from_env", cfg.APIToken)
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("log level: got %q want debug", cfg.LogLevel)
	}
}

func TestFromEnvLoadsDotEnvWhenTokenMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("OPENLANE_API_TOKEN=tola_from_dotenv\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("OPENLANE_ENV_FILE", path)
	t.Setenv("OPENLANE_API_TOKEN", "")
	t.Setenv("OPENLANE_BASE_URL", "")
	t.Setenv("OPENLANE_ORGANIZATION_ID", "")
	t.Setenv("OPENLANE_MCP_LOG_LEVEL", "")
	t.Setenv("OPENLANE_ALLOW_WRITE", "")
	t.Setenv("OPENLANE_ALLOW_DELETE", "")
	t.Setenv("OPENLANE_MCP_TRANSPORT", "")
	t.Setenv("OPENLANE_MCP_HTTP_ADDR", "")

	cfg, err := FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIToken != "tola_from_dotenv" {
		t.Fatalf("token: got %q", cfg.APIToken)
	}
}
