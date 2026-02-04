package logs

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// LogEntry represents a log file with its content
type LogEntry struct {
	Filename string
	Content  string // last 20 lines
}

// FindRecentErrors searches for error log files in the current directory
func FindRecentErrors() []LogEntry {
	var logEntries []LogEntry

	// Common error log filenames to search for
	logFiles := []string{
		"npm-debug.log",
		"yarn-error.log",
		"error.log",
	}

	// Check specific log files
	for _, logFile := range logFiles {
		if _, err := os.Stat(logFile); err == nil {
			content := readLastLines(logFile, 20)
			if content != "" {
				logEntries = append(logEntries, LogEntry{
					Filename: logFile,
					Content:  content,
				})
			}
		}
	}

	// Find all .log files in current directory
	matches, err := filepath.Glob("*.log")
	if err == nil {
		for _, match := range matches {
			// Skip if already processed
			alreadyProcessed := false
			for _, logFile := range logFiles {
				if match == logFile {
					alreadyProcessed = true
					break
				}
			}

			if !alreadyProcessed {
				content := readLastLines(match, 20)
				if content != "" {
					logEntries = append(logEntries, LogEntry{
						Filename: match,
						Content:  content,
					})
				}
			}
		}
	}

	return logEntries
}

// readLastLines reads the last N lines from a file
func readLastLines(filename string, n int) string {
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Get last N lines
	start := 0
	if len(lines) > n {
		start = len(lines) - n
	}

	lastLines := lines[start:]
	return strings.Join(lastLines, "\n")
}
