package project

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// EnvVariable represents an environment variable with its key, value, and secret status
type EnvVariable struct {
	Key      string
	Value    string // either actual value or "[SET]"
	IsSecret bool
}

// ReadEnvFile reads the .env file and returns environment variables with sanitized secrets
func ReadEnvFile() []EnvVariable {
	var envVars []EnvVariable

	file, err := os.Open(".env")
	if err != nil {
		// .env file doesn't exist, return empty slice
		return envVars
	}
	defer file.Close()

	// Regex pattern for environment variables: KEY=VALUE
	envPattern := regexp.MustCompile(`^([A-Z_][A-Z0-9_]*)=(.*)$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Match environment variable pattern
		matches := envPattern.FindStringSubmatch(line)
		if len(matches) != 3 {
			// Malformed line, skip
			continue
		}

		key := matches[1]
		value := matches[2]

		// Check if this is a secret
		isSecret := isSecretKey(key)

		envVar := EnvVariable{
			Key:      key,
			Value:    value,
			IsSecret: isSecret,
		}

		// Sanitize secret values
		if isSecret {
			envVar.Value = "[SET]"
		}

		envVars = append(envVars, envVar)
	}

	return envVars
}

// isSecretKey checks if a key name suggests it contains sensitive information
func isSecretKey(key string) bool {
	keyLower := strings.ToLower(key)

	secretKeywords := []string{
		"key",
		"secret",
		"password",
		"token",
		"database",
		"api",
		"private",
		"auth",
		"credential",
	}

	for _, keyword := range secretKeywords {
		if strings.Contains(keyLower, keyword) {
			return true
		}
	}

	return false
}
