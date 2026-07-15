package checks

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/DSiddharth24/RepoHealth/internal/scanner"
)

const (
	warnSizeBytes = 5 * 1024 * 1024  // 5 MB
	failSizeBytes = 25 * 1024 * 1024 // 25 MB
)

// FileSizeCheck warns or fails for oversized committed files.
type FileSizeCheck struct{}

func (c *FileSizeCheck) Name() string { return "Large files" }

func (c *FileSizeCheck) Run(repoPath string, sc *scanner.Scanner) []Finding {
	var findings []Finding

	err := sc.Walk(func(path string, info fs.FileInfo) error {
		size := info.Size()
		if size <= int64(warnSizeBytes) {
			return nil
		}

		rel, _ := filepath.Rel(repoPath, path)
		if rel == "" {
			rel = path
		}
		relSlash := filepath.ToSlash(rel)

		severity := SeverityWarn
		if size > int64(failSizeBytes) {
			severity = SeverityFail
		}

		findings = append(findings, Finding{
			Check:    c.Name(),
			Severity: severity,
			Message:  fmt.Sprintf("%s (%s)", relSlash, humanSize(size)),
			File:     relSlash,
		})

		return nil
	})

	if err != nil {
		findings = append(findings, Finding{
			Check:    c.Name(),
			Severity: SeverityFail,
			Message:  fmt.Sprintf("Error scanning files: %v", err),
		})
	}

	if len(findings) == 0 {
		findings = append(findings, Finding{
			Check:    c.Name(),
			Severity: SeverityPass,
			Message:  "No oversized files detected",
		})
	}

	return findings
}

// humanSize returns a human-friendly file size string.
func humanSize(bytes int64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
