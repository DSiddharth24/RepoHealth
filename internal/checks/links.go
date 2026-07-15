package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DSiddharth24/RepoHealth/internal/scanner"
)

// LinksCheck scans all .md files for broken internal markdown links.
type LinksCheck struct{}

func (c *LinksCheck) Name() string { return "Broken links" }

// mdLinkRe matches markdown links: [text](target)
// It captures the target in group 1.
var mdLinkRe = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)

// fenceRe matches the opening/closing of a fenced code block line.
var fenceRe = regexp.MustCompile("^\\s*(`{3,}|~{3,})")

// inlineCodeRe matches inline code spans.
var inlineCodeRe = regexp.MustCompile("`[^`]+`")

// stripCodeBlocks removes fenced code blocks and inline code from markdown
// content so that example links inside them are not treated as real links.
func stripCodeBlocks(content string) string {
	var result strings.Builder
	inFence := false
	for _, line := range strings.Split(content, "\n") {
		if fenceRe.MatchString(line) {
			inFence = !inFence
			continue
		}
		if !inFence {
			cleaned := inlineCodeRe.ReplaceAllString(line, "")
			result.WriteString(cleaned)
			result.WriteString("\n")
		}
	}
	return result.String()
}

func (c *LinksCheck) Run(repoPath string, sc *scanner.Scanner) []Finding {
	mdFiles, err := sc.CollectMDFiles()
	if err != nil {
		return []Finding{{
			Check:    c.Name(),
			Severity: SeverityFail,
			Message:  fmt.Sprintf("Error scanning for .md files: %v", err),
		}}
	}

	if len(mdFiles) == 0 {
		return []Finding{{
			Check:    c.Name(),
			Severity: SeverityPass,
			Message:  "No markdown files to check",
		}}
	}

	var broken []Finding

	for _, mdFile := range mdFiles {
		data, err := os.ReadFile(mdFile)
		if err != nil {
			continue // skip files we can't read
		}

		relMdFile, _ := filepath.Rel(repoPath, mdFile)
		if relMdFile == "" {
			relMdFile = mdFile
		}
		mdDir := filepath.Dir(mdFile)

		// Strip code blocks so example links inside them aren't matched
		cleaned := stripCodeBlocks(string(data))
		matches := mdLinkRe.FindAllStringSubmatch(cleaned, -1)
		for _, match := range matches {
			target := match[2]

			// Skip external links
			if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
				continue
			}
			// Skip mailto links
			if strings.HasPrefix(target, "mailto:") {
				continue
			}
			// Skip pure anchor links
			if strings.HasPrefix(target, "#") {
				continue
			}

			// Strip fragment identifier for file existence check
			cleanTarget := target
			if idx := strings.Index(cleanTarget, "#"); idx >= 0 {
				cleanTarget = cleanTarget[:idx]
			}

			// Strip query parameters (rare in local links but be safe)
			if idx := strings.Index(cleanTarget, "?"); idx >= 0 {
				cleanTarget = cleanTarget[:idx]
			}

			if cleanTarget == "" {
				continue // was just a fragment, skip
			}

			// Resolve the target relative to the markdown file's directory
			resolved := filepath.Join(mdDir, filepath.FromSlash(cleanTarget))

			if _, err := os.Stat(resolved); os.IsNotExist(err) {
				broken = append(broken, Finding{
					Check:    c.Name(),
					Severity: SeverityFail,
					Message:  fmt.Sprintf("%s → %s (target does not exist)", filepath.ToSlash(relMdFile), target),
					File:     filepath.ToSlash(relMdFile),
				})
			}
		}
	}

	if len(broken) == 0 {
		return []Finding{{
			Check:    c.Name(),
			Severity: SeverityPass,
			Message:  "All internal links are valid",
		}}
	}

	return broken
}
