package project

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

// Dependencies contains project dependencies
type Dependencies struct {
	Dependencies    map[string]string
	DevDependencies map[string]string
}

// ReadDependencies reads dependencies based on project type
func ReadDependencies(projectType string) Dependencies {
	deps := Dependencies{
		Dependencies:    make(map[string]string),
		DevDependencies: make(map[string]string),
	}

	switch projectType {
	case "Node.js":
		return readNodeDependencies()
	case "Python":
		return readPythonDependencies()
	case "Java":
		return readJavaDependencies()
	case "C#":
		return readCSharpDependencies()
	case "PHP":
		return readPHPDependencies()
	case "Ruby":
		return readRubyDependencies()
	case "Go":
		return readGoDependencies()
	case "Rust":
		return readRustDependencies()
	default:
		return deps
	}
}

// readNodeDependencies reads package.json for Node.js dependencies
func readNodeDependencies() Dependencies {
	deps := Dependencies{
		Dependencies:    make(map[string]string),
		DevDependencies: make(map[string]string),
	}

	data, err := os.ReadFile("package.json")
	if err != nil {
		return deps
	}

	var packageJSON struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return deps
	}

	deps.Dependencies = packageJSON.Dependencies
	deps.DevDependencies = packageJSON.DevDependencies

	return deps
}

// readPythonDependencies reads requirements.txt for Python dependencies
func readPythonDependencies() Dependencies {
	deps := Dependencies{
		Dependencies:    make(map[string]string),
		DevDependencies: make(map[string]string),
	}

	file, err := os.Open("requirements.txt")
	if err != nil {
		return deps
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse package==version or package>=version
		var name, version string
		if strings.Contains(line, "==") {
			parts := strings.Split(line, "==")
			if len(parts) == 2 {
				name = strings.TrimSpace(parts[0])
				version = strings.TrimSpace(parts[1])
			}
		} else if strings.Contains(line, ">=") {
			parts := strings.Split(line, ">=")
			if len(parts) == 2 {
				name = strings.TrimSpace(parts[0])
				version = ">=" + strings.TrimSpace(parts[1])
			}
		} else if strings.Contains(line, "<=") {
			parts := strings.Split(line, "<=")
			if len(parts) == 2 {
				name = strings.TrimSpace(parts[0])
				version = "<=" + strings.TrimSpace(parts[1])
			}
		} else {
			// No version specified
			name = line
			version = "*"
		}

		if name != "" {
			deps.Dependencies[name] = version
		}
	}

	return deps
}

// readGoDependencies reads go.mod for Go dependencies
func readGoDependencies() Dependencies {
	deps := Dependencies{
		Dependencies:    make(map[string]string),
		DevDependencies: make(map[string]string),
	}

	file, err := os.Open("go.mod")
	if err != nil {
		return deps
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inRequireBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Detect require block
		if strings.HasPrefix(line, "require (") {
			inRequireBlock = true
			continue
		}

		// End of require block
		if inRequireBlock && line == ")" {
			inRequireBlock = false
			continue
		}

		// Parse require lines
		if inRequireBlock || strings.HasPrefix(line, "require ") {
			// Remove "require " prefix if present
			line = strings.TrimPrefix(line, "require ")

			// Split by whitespace
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[0]
				version := parts[1]

				// Skip indirect dependencies in summary (we'll include them but mark them)
				deps.Dependencies[name] = version
			}
		}
	}

	return deps
}

// readRustDependencies reads Cargo.toml for Rust dependencies
func readRustDependencies() Dependencies {
	deps := Dependencies{
		Dependencies:    make(map[string]string),
		DevDependencies: make(map[string]string),
	}

	file, err := os.Open("Cargo.toml")
	if err != nil {
		return deps
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inDependencies := false
	inDevDependencies := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Detect sections
		if line == "[dependencies]" {
			inDependencies = true
			inDevDependencies = false
			continue
		}

		if line == "[dev-dependencies]" {
			inDependencies = false
			inDevDependencies = true
			continue
		}

		// Exit section if we hit another section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			inDependencies = false
			inDevDependencies = false
			continue
		}

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse dependency lines
		if inDependencies || inDevDependencies {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				version := strings.TrimSpace(parts[1])
				// Remove quotes
				version = strings.Trim(version, "\"")

				if inDependencies {
					deps.Dependencies[name] = version
				} else {
					deps.DevDependencies[name] = version
				}
			}
		}
	}

	return deps
}
