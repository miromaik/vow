package git

import (
	"os/exec"
	"strings"
)

// GitInfo contains Git repository information
type GitInfo struct {
	IsRepo         bool
	Branch         string
	LastCommit     string
	LastCommitTime string
	ModifiedFiles  int
	RemoteURL      string
}

// GetGitInfo collects Git repository information
func GetGitInfo() GitInfo {
	info := GitInfo{
		IsRepo: false,
	}

	// Check if we're in a git repository
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		// Not a git repository
		return info
	}

	info.IsRepo = true

	// Get current branch
	info.Branch = getCurrentBranch()

	// Get last commit information
	info.LastCommit, info.LastCommitTime = getLastCommit()

	// Get modified files count
	info.ModifiedFiles = getModifiedFilesCount()

	// Get remote URL
	info.RemoteURL = getRemoteURL()

	return info
}

// getCurrentBranch returns the current Git branch name
func getCurrentBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

// getLastCommit returns the last commit message and relative time
func getLastCommit() (string, string) {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%s|||%ar")
	output, err := cmd.Output()
	if err != nil {
		return "No commits", ""
	}

	parts := strings.Split(string(output), "|||")
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}

	return strings.TrimSpace(string(output)), ""
}

// getModifiedFilesCount returns the number of modified files
func getModifiedFilesCount() int {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return 0
	}

	return len(lines)
}

// getRemoteURL returns the remote origin URL with credentials sanitized
func getRemoteURL() string {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "No remote"
	}

	url := strings.TrimSpace(string(output))

	// Sanitize credentials from URL
	url = sanitizeGitURL(url)

	return url
}

// sanitizeGitURL removes credentials from Git URLs
func sanitizeGitURL(url string) string {
	// Handle HTTPS URLs with credentials: https://user:pass@github.com/...
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		// Find @ symbol
		atIndex := strings.Index(url, "@")
		if atIndex != -1 {
			// Find protocol end
			protocolEnd := strings.Index(url, "://")
			if protocolEnd != -1 && protocolEnd < atIndex {
				// Reconstruct URL without credentials
				protocol := url[:protocolEnd+3]
				remainder := url[atIndex+1:]
				return protocol + remainder
			}
		}
	}

	return url
}
