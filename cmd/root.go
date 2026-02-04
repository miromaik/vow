package cmd

import (
	"fmt"
	"os"
	"time"

	"vow/internal/git"
	"vow/internal/logs"
	"vow/internal/output"
	"vow/internal/project"
	"vow/internal/system"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	outputFile string
	noLogs     bool
	format     string
)

var rootCmd = &cobra.Command{
	Use:   "vow",
	Short: "VOW - Version, OS, Workspace report generator",
	Long: `VOW automatically collects debugging information about your project and system,
then generates a clean markdown report ready to share on forums, GitHub issues, or Discord.`,
	Run: generateReport,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "REPORT.txt", "Output filename for the report")
	rootCmd.Flags().BoolVar(&noLogs, "no-logs", false, "Skip error logs section")
	rootCmd.Flags().StringVarP(&format, "format", "f", "txt", "Output format: txt or markdown")
}

func generateReport(cmd *cobra.Command, args []string) {
	color.Cyan("Collecting debugging information...")

	// Collect system information
	systemInfo := system.GetSystemInfo()
	color.Green("  System information collected")

	// Detect project type
	projectType := project.DetectProjectType()
	if projectType.Type != "Unknown" {
		color.Green("  Project type detected: %s", projectType.Type)
	} else {
		color.Yellow("  No project type detected")
	}

	// Collect environment variables
	envVars := project.ReadEnvFile()
	if len(envVars) > 0 {
		secretCount := 0
		for _, env := range envVars {
			if env.IsSecret {
				secretCount++
			}
		}
		color.Green("  Environment variables read: %d (%d secrets hidden)", len(envVars), secretCount)
	}

	// Collect dependencies
	dependencies := project.ReadDependencies(projectType.Type)
	totalDeps := len(dependencies.Dependencies) + len(dependencies.DevDependencies)
	if totalDeps > 0 {
		color.Green("  Dependencies collected: %d", totalDeps)
	}

	// Collect git information
	gitInfo := git.GetGitInfo()
	if gitInfo.IsRepo {
		color.Green("  Git information collected")
	}

	// Collect recent errors
	var logEntries []logs.LogEntry
	if !noLogs {
		logEntries = logs.FindRecentErrors()
		if len(logEntries) > 0 {
			color.Green("  Error logs found: %d files", len(logEntries))
		}
	}

	// Generate markdown report
	reportData := output.ReportData{
		Timestamp:    time.Now().UTC().Format("2006-01-02 15:04:05"),
		System:       systemInfo,
		Project:      projectType,
		EnvVars:      envVars,
		Dependencies: dependencies,
		Git:          gitInfo,
		Logs:         logEntries,
	}

	// Generate report based on format
	var reportContent string
	if format == "markdown" || format == "md" {
		reportContent = output.GenerateMarkdown(reportData)
	} else {
		reportContent = output.GenerateText(reportData)
	}

	// Save to file
	err := os.WriteFile(outputFile, []byte(reportContent), 0644)
	if err != nil {
		color.Red("Error writing report: %v", err)
		os.Exit(1)
	}

	// Print success message
	fmt.Println()
	color.Green("Report generated successfully: %s", outputFile)
	fmt.Println()
	color.Cyan("Report includes:")

	fmt.Printf("  - System information\n")
	if projectType.Type != "Unknown" {
		fmt.Printf("  - Project type: %s", projectType.Type)
		if projectType.Version != "" {
			fmt.Printf(" %s", projectType.Version)
		}
		fmt.Println()
	}
	if totalDeps > 0 {
		fmt.Printf("  - %d dependencies\n", totalDeps)
	}
	if len(envVars) > 0 {
		secretCount := 0
		for _, env := range envVars {
			if env.IsSecret {
				secretCount++
			}
		}
		fmt.Printf("  - %d environment variables", len(envVars))
		if secretCount > 0 {
			fmt.Printf(" (%d secrets hidden)", secretCount)
		}
		fmt.Println()
	}
	if gitInfo.IsRepo {
		fmt.Printf("  - Git branch: %s\n", gitInfo.Branch)
	}
	if len(logEntries) > 0 {
		fmt.Printf("  - %d log files analyzed\n", len(logEntries))
	}

	fmt.Println()
	color.Green("Ready to share on StackOverflow, GitHub Issues, or Discord!")
}
