package project

import (
	"os"
	"os/exec"
	"strings"
)

// ProjectType contains information about the detected project type
type ProjectType struct {
	Type           string // "Node.js", "Python", "Go", "Rust", "Unknown"
	Version        string // e.g., "v20.10.0"
	PackageManager string // "npm", "yarn", "pnpm", "pip", etc.
	PMVersion      string
}

// DetectProjectType detects the type of project in the current directory
func DetectProjectType() ProjectType {
	project := ProjectType{
		Type: "Unknown",
	}

	// Check for Node.js
	if fileExists("package.json") {
		project.Type = "Node.js"
		project.Version = getNodeVersion()
		project.PackageManager, project.PMVersion = getNodePackageManager()
		return project
	}

	// Check for Python
	if fileExists("requirements.txt") || fileExists("pyproject.toml") {
		project.Type = "Python"
		project.Version = getPythonVersion()
		project.PackageManager = "pip"
		project.PMVersion = getPipVersion()
		return project
	}

	// Check for Java (Maven)
	if fileExists("pom.xml") {
		project.Type = "Java"
		project.Version = getJavaVersion()
		project.PackageManager = "maven"
		project.PMVersion = getMavenVersion()
		return project
	}

	// Check for Java (Gradle)
	if fileExists("build.gradle") || fileExists("build.gradle.kts") {
		project.Type = "Java"
		project.Version = getJavaVersion()
		project.PackageManager = "gradle"
		project.PMVersion = getGradleVersion()
		return project
	}

	// Check for C# (.NET)
	if hasCSharpProject() {
		project.Type = "C#"
		project.Version = getDotnetVersion()
		project.PackageManager = "dotnet"
		return project
	}

	// Check for PHP
	if fileExists("composer.json") {
		project.Type = "PHP"
		project.Version = getPHPVersion()
		project.PackageManager = "composer"
		project.PMVersion = getComposerVersion()
		return project
	}

	// Check for Ruby
	if fileExists("Gemfile") {
		project.Type = "Ruby"
		project.Version = getRubyVersion()
		project.PackageManager = "bundler"
		project.PMVersion = getBundlerVersion()
		return project
	}

	// Check for Go
	if fileExists("go.mod") {
		project.Type = "Go"
		project.Version = getGoVersion()
		project.PackageManager = "go"
		return project
	}

	// Check for Rust
	if fileExists("Cargo.toml") {
		project.Type = "Rust"
		project.Version = getRustVersion()
		project.PackageManager = "cargo"
		project.PMVersion = getCargoVersion()
		return project
	}

	return project
}

// fileExists checks if a file exists in the current directory
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// getNodeVersion returns the Node.js version
func getNodeVersion() string {
	cmd := exec.Command("node", "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// getNodePackageManager detects and returns the package manager and version
func getNodePackageManager() (string, string) {
	// Check for pnpm
	if fileExists("pnpm-lock.yaml") {
		cmd := exec.Command("pnpm", "--version")
		if output, err := cmd.Output(); err == nil {
			return "pnpm", strings.TrimSpace(string(output))
		}
	}

	// Check for yarn
	if fileExists("yarn.lock") {
		cmd := exec.Command("yarn", "--version")
		if output, err := cmd.Output(); err == nil {
			return "yarn", strings.TrimSpace(string(output))
		}
	}

	// Default to npm
	cmd := exec.Command("npm", "--version")
	if output, err := cmd.Output(); err == nil {
		return "npm", strings.TrimSpace(string(output))
	}

	return "npm", ""
}

// getPythonVersion returns the Python version
func getPythonVersion() string {
	// Try python3 first
	cmd := exec.Command("python3", "--version")
	output, err := cmd.Output()
	if err == nil {
		version := strings.TrimSpace(string(output))
		// Python outputs "Python 3.x.x", so we trim the "Python " prefix
		return strings.TrimPrefix(version, "Python ")
	}

	// Try python
	cmd = exec.Command("python", "--version")
	output, err = cmd.Output()
	if err == nil {
		version := strings.TrimSpace(string(output))
		return strings.TrimPrefix(version, "Python ")
	}

	return ""
}

// getPipVersion returns the pip version
func getPipVersion() string {
	cmd := exec.Command("pip", "--version")
	output, err := cmd.Output()
	if err != nil {
		cmd = exec.Command("pip3", "--version")
		output, err = cmd.Output()
		if err != nil {
			return ""
		}
	}

	// pip outputs "pip X.Y.Z from ...", extract just the version
	parts := strings.Fields(string(output))
	if len(parts) >= 2 {
		return parts[1]
	}

	return ""
}

// getGoVersion returns the Go version
func getGoVersion() string {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// go version outputs "go version goX.Y.Z OS/ARCH"
	parts := strings.Fields(string(output))
	if len(parts) >= 3 {
		return parts[2]
	}

	return ""
}

// getRustVersion returns the Rust version
func getRustVersion() string {
	cmd := exec.Command("rustc", "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// rustc outputs "rustc X.Y.Z (hash date)"
	parts := strings.Fields(string(output))
	if len(parts) >= 2 {
		return parts[1]
	}

	return ""
}

// getCargoVersion returns the Cargo version
func getCargoVersion() string {
	cmd := exec.Command("cargo", "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// cargo outputs "cargo X.Y.Z (hash date)"
	parts := strings.Fields(string(output))
	if len(parts) >= 2 {
		return parts[1]
	}

	return ""
}

// hasCSharpProject checks if there's a C# project file
func hasCSharpProject() bool {
	files, err := os.ReadDir(".")
	if err != nil {
		return false
	}

	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if strings.HasSuffix(name, ".csproj") || strings.HasSuffix(name, ".sln") {
				return true
			}
		}
	}
	return false
}

// getJavaVersion returns the Java version
func getJavaVersion() string {
	cmd := exec.Command("java", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	// Java outputs version on stderr: 'java version "17.0.1"' or 'openjdk version "11.0.12"'
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "version") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	}

	return ""
}

// getMavenVersion returns the Maven version
func getMavenVersion() string {
	cmd := exec.Command("mvn", "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// Maven outputs "Apache Maven X.Y.Z ..."
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 3 {
			return parts[2]
		}
	}

	return ""
}

// getGradleVersion returns the Gradle version
func getGradleVersion() string {
	cmd := exec.Command("gradle", "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// Gradle outputs multiple lines, version is in "Gradle X.Y"
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Gradle ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	}

	return ""
}

// getDotnetVersion returns the .NET version
func getDotnetVersion() string {
	cmd := exec.Command("dotnet", "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

// getPHPVersion returns the PHP version
func getPHPVersion() string {
	cmd := exec.Command("php", "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// PHP outputs "PHP X.Y.Z ..."
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 {
			return parts[1]
		}
	}

	return ""
}

// getComposerVersion returns the Composer version
func getComposerVersion() string {
	cmd := exec.Command("composer", "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// Composer outputs "Composer version X.Y.Z ..."
	parts := strings.Fields(string(output))
	if len(parts) >= 3 {
		return parts[2]
	}

	return ""
}

// getRubyVersion returns the Ruby version
func getRubyVersion() string {
	cmd := exec.Command("ruby", "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// Ruby outputs "ruby X.Y.Z ..."
	parts := strings.Fields(string(output))
	if len(parts) >= 2 {
		return parts[1]
	}

	return ""
}

// getBundlerVersion returns the Bundler version
func getBundlerVersion() string {
	cmd := exec.Command("bundle", "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// Bundler outputs "Bundler version X.Y.Z"
	parts := strings.Fields(string(output))
	if len(parts) >= 3 {
		return parts[2]
	}

	return ""
}
