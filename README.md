# VOW - Version, OS, Workspace Report

A command-line tool that automatically collects debugging information about your project and system, then generates a report ready to share on forums or GitHub Issues.

## Features

- Detects project type (Node.js, Python, Java, C#, PHP, Ruby, Go, Rust)
- Collects system info, dependencies, environment variables, git status
- Hides sensitive data automatically (API keys, passwords, etc.)
- Output formats: plain text (default) or markdown
- Cross-platform (Linux, macOS, Windows)

## Installation

### Quick Install (macOS/Linux)

```bash
curl -sSL https://raw.githubusercontent.com/miromaik/vow/main/install.sh | bash
```

### Using Go

```bash
go install github.com/miromaik/vow@latest
```

Make sure `$GOPATH/bin` is in your PATH:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Manual Build

```bash
git clone https://github.com/miromaik/vow.git
cd vow
make install
```

## Usage

```bash
# Basic usage (generates REPORT.txt)
vow

# Custom output file
vow --output my-report.txt
vow -o debug.txt

# Generate markdown instead of plain text
vow --format markdown -o report.md
vow -f md

# Skip error logs
vow --no-logs
```

## What It Collects

- **System**: OS, version, architecture, shell, hostname
- **Project**: Type, version, package manager, dependencies
- **Environment**: Variables from .env (secrets hidden)
- **Git**: Branch, last commit, modified files, remote URL
- **Logs**: Recent error logs (npm-debug.log, *.log files)

## Example Output (Plain Text)

```
================================================================================
VOW DEBUG REPORT
================================================================================
Generated: 2026-02-04 14:30:00 UTC

SYSTEM INFORMATION
--------------------------------------------------------------------------------
Operating System: darwin 14.2.1
Architecture:     arm64
Shell:            /bin/zsh
Hostname:         developer-macbook.local

PROJECT INFORMATION
--------------------------------------------------------------------------------
Type: Go
Version: go1.21.5

Dependencies:
  - github.com/spf13/cobra: v1.8.0
  - github.com/fatih/color: v1.16.0

ENVIRONMENT VARIABLES
--------------------------------------------------------------------------------
NODE_ENV: production
API_KEY: [SET] (hidden)
DATABASE_URL: [SET] (hidden)

GIT INFORMATION
--------------------------------------------------------------------------------
Branch: main
Last Commit: Add new feature (2 hours ago)
Modified Files: 3
Remote: https://github.com/user/project.git
```

## Privacy

- Automatically hides values containing: KEY, SECRET, PASSWORD, TOKEN, DATABASE, API, PRIVATE
- Removes credentials from Git URLs
- Everything runs locally, no data sent anywhere

## Requirements

- Go 1.21+ (for building)
- Git (optional, for git info)
- Project tools (node, python, go, rustc) for version detection

## License

MIT License - see LICENSE file

## Contributing

PRs welcome! To add support for new languages, update:
- `internal/project/detector.go` - detection logic
- `internal/project/deps.go` - dependency parsing

