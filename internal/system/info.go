package system

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// SystemInfo contains system-level information
type SystemInfo struct {
	OS           string
	OSVersion    string
	Architecture string
	Shell        string
	Hostname     string
}

// GetSystemInfo collects system information (OS, architecture, shell, hostname)
func GetSystemInfo() SystemInfo {
	info := SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Shell:        getShell(),
		Hostname:     getHostname(),
		OSVersion:    getOSVersion(),
	}

	return info
}

// getShell returns the current shell from environment variable
func getShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "unknown"
	}
	return shell
}

// getHostname returns the system hostname
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// getOSVersion returns the OS version based on the platform
func getOSVersion() string {
	switch runtime.GOOS {
	case "linux":
		return getLinuxVersion()
	case "darwin":
		return getMacOSVersion()
	case "windows":
		return getWindowsVersion()
	default:
		return "Unknown version"
	}
}

// getLinuxVersion reads /etc/os-release to get the Linux distribution version
func getLinuxVersion() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "Unknown version"
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			version := strings.TrimPrefix(line, "PRETTY_NAME=")
			version = strings.Trim(version, "\"")
			return version
		}
	}

	return "Unknown version"
}

// getMacOSVersion executes sw_vers to get macOS version
func getMacOSVersion() string {
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err != nil {
		return "Unknown version"
	}

	return strings.TrimSpace(string(output))
}

// getWindowsVersion executes ver command to get Windows version
func getWindowsVersion() string {
	cmd := exec.Command("cmd", "/c", "ver")
	output, err := cmd.Output()
	if err != nil {
		return "Unknown version"
	}

	return strings.TrimSpace(string(output))
}
