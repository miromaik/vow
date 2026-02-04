package project

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

// readJavaDependencies reads pom.xml or build.gradle for Java dependencies
func readJavaDependencies() Dependencies {
	deps := Dependencies{
		Dependencies:    make(map[string]string),
		DevDependencies: make(map[string]string),
	}

	// Try Maven first
	if data, err := os.ReadFile("pom.xml"); err == nil {
		content := string(data)
		lines := strings.Split(content, "\n")

		var artifactId, version string
		for _, line := range lines {
			line = strings.TrimSpace(line)

			if strings.Contains(line, "<artifactId>") && !strings.Contains(line, "${") {
				artifactId = extractXMLValue(line, "artifactId")
			}
			if strings.Contains(line, "<version>") && !strings.Contains(line, "${") {
				version = extractXMLValue(line, "version")
			}
			if strings.Contains(line, "</dependency>") && artifactId != "" {
				deps.Dependencies[artifactId] = version
				artifactId = ""
				version = ""
			}
		}
	}

	// Try Gradle
	if len(deps.Dependencies) == 0 {
		for _, filename := range []string{"build.gradle", "build.gradle.kts"} {
			if data, err := os.ReadFile(filename); err == nil {
				content := string(data)
				lines := strings.Split(content, "\n")

				for _, line := range lines {
					line = strings.TrimSpace(line)
					if strings.Contains(line, "implementation") || strings.Contains(line, "compile") {
						parts := strings.Split(line, "'")
						if len(parts) >= 2 {
							depParts := strings.Split(parts[1], ":")
							if len(depParts) >= 2 {
								name := depParts[len(depParts)-2]
								version := "*"
								if len(depParts) >= 3 {
									version = depParts[len(depParts)-1]
								}
								deps.Dependencies[name] = version
							}
						}
					}
				}
				break
			}
		}
	}

	return deps
}

// readCSharpDependencies reads .csproj for C# dependencies
func readCSharpDependencies() Dependencies {
	deps := Dependencies{
		Dependencies: make(map[string]string),
	}

	files, err := os.ReadDir(".")
	if err != nil {
		return deps
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csproj") {
			data, err := os.ReadFile(file.Name())
			if err != nil {
				continue
			}

			content := string(data)
			lines := strings.Split(content, "\n")

			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.Contains(line, "PackageReference") {
					name := extractAttribute(line, "Include")
					version := extractAttribute(line, "Version")
					if name != "" {
						deps.Dependencies[name] = version
					}
				}
			}
			break
		}
	}

	return deps
}

// readPHPDependencies reads composer.json for PHP dependencies
func readPHPDependencies() Dependencies {
	deps := Dependencies{
		Dependencies:    make(map[string]string),
		DevDependencies: make(map[string]string),
	}

	data, err := os.ReadFile("composer.json")
	if err != nil {
		return deps
	}

	var composerJSON struct {
		Require    map[string]string `json:"require"`
		RequireDev map[string]string `json:"require-dev"`
	}

	if err := json.Unmarshal(data, &composerJSON); err != nil {
		return deps
	}

	deps.Dependencies = composerJSON.Require
	deps.DevDependencies = composerJSON.RequireDev

	return deps
}

// readRubyDependencies reads Gemfile for Ruby dependencies
func readRubyDependencies() Dependencies {
	deps := Dependencies{
		Dependencies:    make(map[string]string),
		DevDependencies: make(map[string]string),
	}

	file, err := os.Open("Gemfile")
	if err != nil {
		return deps
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "gem ") {
			line = strings.TrimPrefix(line, "gem ")
			line = strings.TrimSpace(line)

			parts := strings.Split(line, ",")
			if len(parts) >= 1 {
				name := strings.Trim(parts[0], "\"' ")
				version := "*"

				if len(parts) >= 2 {
					versionPart := strings.TrimSpace(parts[1])
					versionPart = strings.Trim(versionPart, "\"' ")
					if versionPart != "" {
						version = versionPart
					}
				}

				deps.Dependencies[name] = version
			}
		}
	}

	return deps
}

// extractXMLValue extracts value from XML tag
func extractXMLValue(line, tag string) string {
	start := strings.Index(line, "<"+tag+">")
	end := strings.Index(line, "</"+tag+">")
	if start != -1 && end != -1 {
		return line[start+len(tag)+2 : end]
	}
	return ""
}

// extractAttribute extracts attribute value from XML/HTML-like string
func extractAttribute(line, attr string) string {
	pattern := attr + "=\""
	start := strings.Index(line, pattern)
	if start == -1 {
		return ""
	}
	start += len(pattern)
	end := strings.Index(line[start:], "\"")
	if end == -1 {
		return ""
	}
	return line[start : start+end]
}
