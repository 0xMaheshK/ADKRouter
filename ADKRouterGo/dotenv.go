package adkrouter

import (
	"bufio"
	"os"
	"strings"
)

// loadDotEnv loads a .env file into the process environment following the
// python-dotenv conventions: KEY=VALUE lines, optional `export` prefix,
// comments starting with `#`, and quoted values. It mirrors the optional
// dotenv support in the Python implementation without introducing an
// external dependency. A no-op when envPath does not exist.
func loadDotEnv(envPath string) error {
	f, err := os.Open(envPath)
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
		line = strings.TrimPrefix(line, "export ")
		line = strings.TrimSpace(line)

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
	return scanner.Err()
}

// LoadEnvConfig loads a .env file into the process environment. It is a
// no-op when the file does not exist. Existing environment variables are
// preserved (never overridden), matching python-dotenv's default behavior.
func LoadEnvConfig(envPath string) error {
	if envPath == "" {
		envPath = ".env"
	}
	return loadDotEnv(envPath)
}
