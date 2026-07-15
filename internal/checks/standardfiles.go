package checks

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/DSiddharth24/RepoHealth/internal/scanner"
)

// StandardFilesCheck verifies the presence of LICENSE and .gitignore.
type StandardFilesCheck struct{}

func (c *StandardFilesCheck) Name() string { return "Standard files" }

func (c *StandardFilesCheck) Run(repoPath string, sc *scanner.Scanner) []Finding {
	hasLicense := fileExistsCaseInsensitive(repoPath, "LICENSE") ||
		fileExistsCaseInsensitive(repoPath, "LICENSE.md") ||
		fileExistsCaseInsensitive(repoPath, "LICENSE.txt")

	hasGitignore := fileExists(filepath.Join(repoPath, ".gitignore"))

	var findings []Finding

	if !hasLicense {
		findings = append(findings, Finding{
			Check:    c.Name(),
			Severity: SeverityWarn,
			Message:  "No LICENSE file found",
		})
	}

	if !hasGitignore {
		findings = append(findings, Finding{
			Check:    c.Name(),
			Severity: SeverityWarn,
			Message:  "No .gitignore file found",
		})
	}

	if len(findings) == 0 {
		findings = append(findings, Finding{
			Check:    c.Name(),
			Severity: SeverityPass,
			Message:  "LICENSE and .gitignore present",
		})
	}

	return findings
}

// fileExists checks if a file exists at the given path.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// fileExistsCaseInsensitive checks for a file matching the name case-insensitively
// in the given directory. Only checks immediate children, not subdirectories.
func fileExistsCaseInsensitive(dir, name string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.EqualFold(e.Name(), name) {
			return true
		}
	}
	return false
}
