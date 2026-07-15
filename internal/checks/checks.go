package checks

import (
	"github.com/DSiddharth24/RepoHealth/internal/scanner"
)

// Severity levels for findings.
const (
	SeverityPass = "pass"
	SeverityWarn = "warn"
	SeverityFail = "fail"
)

// Finding represents a single result from a health check.
type Finding struct {
	Check    string `json:"check"`
	Severity string `json:"severity"` // "pass", "warn", or "fail"
	Message  string `json:"message"`
	File     string `json:"file,omitempty"` // optional, empty if repo-wide
}

// Check is the interface that all health checks implement.
// Adding a new check is a matter of implementing this interface
// and registering it in Registry().
type Check interface {
	// Name returns a human-readable name for the check (e.g. "README").
	Name() string
	// Run executes the check against the given repo path and returns findings.
	// The scanner is provided for checks that need to walk the file tree.
	Run(repoPath string, sc *scanner.Scanner) []Finding
}

// Registry returns all registered checks in their preferred display order.
// To add a new check, implement the Check interface and append it here.
func Registry() []Check {
	return []Check{
		&ReadmeCheck{},
		&LinksCheck{},
		&FileSizeCheck{},
		&StandardFilesCheck{},
	}
}
