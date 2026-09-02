package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// loadDotEnvFiles loads KEY=VALUE pairs from the first existing .env file.
// Existing process environment variables are never overwritten.
func loadDotEnvFiles() {
	for _, path := range dotEnvCandidates() {
		if err := loadDotEnvFile(path); err == nil {
			return
		}
	}
}

func dotEnvCandidates() []string {
	var paths []string
	if custom := strings.TrimSpace(os.Getenv("OPENLANE_ENV_FILE")); custom != "" {
		paths = append(paths, custom)
	}
	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(cwd, ".env"))
	}
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		paths = append(paths,
			filepath.Join(exeDir, ".env"),
			filepath.Join(exeDir, "..", ".env"),
		)
	}
	return paths
}

func loadDotEnvFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" || os.Getenv(key) != "" {
			continue
		}
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		_ = os.Setenv(key, value)
	}
	return scanner.Err()
}
