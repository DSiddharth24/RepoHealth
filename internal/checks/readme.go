package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DSiddharth24/RepoHealth/internal/scanner"
)

// ReadmeCheck verifies the presence and quality of README.md.
type ReadmeCheck struct{}

func (c *ReadmeCheck) Name() string { return "README" }

func (c *ReadmeCheck) Run(repoPath string, sc *scanner.Scanner) []Finding {
	readmePath := filepath.Join(repoPath, "README.md")

	data, err := os.ReadFile(readmePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Finding{{
				Check:    c.Name(),
				Severity: SeverityFail,
				Message:  "README.md is missing",
			}}
		}
		return []Finding{{
			Check:    c.Name(),
			Severity: SeverityFail,
			Message:  fmt.Sprintf("Could not read README.md: %v", err),
			File:     "README.md",
		}}
	}

	content := string(data)
	words := countWords(content)
	headings := countHeadings(content)

	var findings []Finding

	if words < 150 {
		findings = append(findings, Finding{
			Check:    c.Name(),
			Severity: SeverityWarn,
			Message:  fmt.Sprintf("README.md is likely a stub (%d words, recommend ≥150)", words),
			File:     "README.md",
		})
	}

	if headings == 0 {
		findings = append(findings, Finding{
			Check:    c.Name(),
			Severity: SeverityWarn,
			Message:  "README.md has no headings (likely unstructured)",
			File:     "README.md",
		})
	}

	if len(findings) == 0 {
		findings = append(findings, Finding{
			Check:    c.Name(),
			Severity: SeverityPass,
			Message:  fmt.Sprintf("Present, well-structured (%d words, %d headings)", words, headings),
			File:     "README.md",
		})
	}

	return findings
}

// countWords counts whitespace-delimited words in text.
func countWords(text string) int {
	fields := strings.Fields(text)
	return len(fields)
}

// headingRe matches markdown headings (lines starting with one or more #).
var headingRe = regexp.MustCompile(`(?m)^#{1,6}\s+\S`)

// countHeadings counts markdown headings in text.
func countHeadings(text string) int {
	return len(headingRe.FindAllString(text, -1))
}
