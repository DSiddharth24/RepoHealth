package cmd

import (
	"fmt"
	"os"

	"github.com/DSiddharth24/RepoHealth/internal/checks"
	"github.com/DSiddharth24/RepoHealth/internal/report"
	"github.com/DSiddharth24/RepoHealth/internal/scanner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	strict     bool
	ignoreList []string
	noColor    bool
)

var rootCmd = &cobra.Command{
	Use:   "repohealth [path]",
	Short: "Scan a repository and report on its health",
	Long: `RepoHealth scans a local repository folder and reports on its "health" —
README quality, broken internal markdown links, oversized committed files,
and missing standard files (LICENSE, .gitignore).

Point it at any directory and get an honest, fast report on whether it
looks like a well-maintained repo.`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE:         runScan,
}

func init() {
	rootCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output machine-readable JSON instead of terminal report")
	rootCmd.Flags().BoolVar(&strict, "strict", false, "Treat all warnings as failures (for CI gating)")
	rootCmd.Flags().StringSliceVar(&ignoreList, "ignore", nil, "Additional glob patterns to skip (repeatable)")
	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable ANSI colors (for log files/CI output)")
}

// Execute runs the root command and returns the exit code.
// We don't call os.Exit here — main.go does that.
func Execute() int {
	exitCode := 0
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		code, err := runScanInner(args)
		exitCode = code
		return err
	}

	if err := rootCmd.Execute(); err != nil {
		if exitCode == 0 {
			exitCode = 2 // tool error
		}
	}

	return exitCode
}

func runScan(cmd *cobra.Command, args []string) error {
	_, err := runScanInner(args)
	return err
}

func runScanInner(args []string) (int, error) {
	if noColor {
		color.NoColor = true
	}

	// Determine repo path
	repoPath := "."
	if len(args) > 0 {
		repoPath = args[0]
	}

	// Validate path
	info, err := os.Stat(repoPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 2, fmt.Errorf("path does not exist: %s", repoPath)
		}
		if os.IsPermission(err) {
			return 2, fmt.Errorf("permission denied: %s", repoPath)
		}
		return 2, fmt.Errorf("cannot access path: %w", err)
	}
	if !info.IsDir() {
		return 2, fmt.Errorf("path is not a directory: %s", repoPath)
	}

	// Create scanner
	sc, err := scanner.NewScanner(repoPath, ignoreList)
	if err != nil {
		return 2, fmt.Errorf("scanner initialization failed: %w", err)
	}

	// Print scanner warnings (e.g. malformed .gitignore)
	for _, w := range sc.Warnings {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", w)
	}

	// Run all registered checks
	registry := checks.Registry()
	var allFindings []checks.Finding

	for _, chk := range registry {
		results := chk.Run(repoPath, sc)
		allFindings = append(allFindings, results...)
	}

	// Determine exit code
	exitCode := determineExitCode(allFindings, strict)

	// Output
	if jsonOutput {
		if err := report.JSON(os.Stdout, allFindings, exitCode); err != nil {
			return 2, fmt.Errorf("failed to write JSON report: %w", err)
		}
	} else {
		report.Terminal(os.Stdout, allFindings, repoPath)
		// Print exit code hint for CI
		if exitCode != 0 {
			fmt.Fprintf(os.Stderr, "Exit code: %d\n", exitCode)
		}
	}

	return exitCode, nil
}

// determineExitCode returns 0 if all checks pass (warnings allowed unless strict),
// 1 if any check fails (or any warning with strict mode).
func determineExitCode(findings []checks.Finding, strict bool) int {
	for _, f := range findings {
		if f.Severity == checks.SeverityFail {
			return 1
		}
		if strict && f.Severity == checks.SeverityWarn {
			return 1
		}
	}
	return 0
}
