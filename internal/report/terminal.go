package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/DSiddharth24/RepoHealth/internal/checks"
	"github.com/fatih/color"
)

// Terminal writes a color-coded, grouped terminal report to the given writer.
func Terminal(w io.Writer, findings []checks.Finding, repoPath string) {
	fmt.Fprintf(w, "\n%s\n\n",
		color.New(color.Bold).Sprintf("RepoHealth — scanning %s", repoPath))

	// Group findings by check name, preserving insertion order.
	type group struct {
		name     string
		findings []checks.Finding
	}
	var groups []group
	seen := map[string]int{}

	for _, f := range findings {
		if idx, ok := seen[f.Check]; ok {
			groups[idx].findings = append(groups[idx].findings, f)
		} else {
			seen[f.Check] = len(groups)
			groups = append(groups, group{name: f.Check, findings: []checks.Finding{f}})
		}
	}

	passed, warnings, failed := 0, 0, 0

	greenCheck := color.New(color.FgGreen, color.Bold)
	yellowWarn := color.New(color.FgYellow, color.Bold)
	redFail := color.New(color.FgRed, color.Bold)
	dim := color.New(color.Faint)

	for _, g := range groups {
		worst := worstSeverity(g.findings)

		switch worst {
		case checks.SeverityPass:
			passed++
			msg := g.findings[0].Message
			fmt.Fprintf(w, "  %s %-24s %s\n",
				greenCheck.Sprint("✓"),
				color.New(color.Bold).Sprint(g.name),
				msg)

		case checks.SeverityWarn:
			warnings++
			warnCount := countSeverity(g.findings, checks.SeverityWarn)
			fmt.Fprintf(w, "  %s %-24s %d warning(s)\n",
				yellowWarn.Sprint("⚠"),
				color.New(color.Bold).Sprint(g.name),
				warnCount)
			for _, f := range g.findings {
				if f.Severity == checks.SeverityWarn {
					fmt.Fprintf(w, "      %s\n", dim.Sprint(f.Message))
				}
			}

		case checks.SeverityFail:
			failed++
			failCount := countSeverity(g.findings, checks.SeverityFail)
			warnCount := countSeverity(g.findings, checks.SeverityWarn)
			parts := []string{fmt.Sprintf("%d issue(s) found", failCount)}
			if warnCount > 0 {
				parts = append(parts, fmt.Sprintf("%d warning(s)", warnCount))
			}
			fmt.Fprintf(w, "  %s %-24s %s\n",
				redFail.Sprint("✗"),
				color.New(color.Bold).Sprint(g.name),
				strings.Join(parts, ", "))
			for _, f := range g.findings {
				if f.Severity == checks.SeverityFail {
					fmt.Fprintf(w, "      %s\n", dim.Sprint(f.Message))
				}
			}
			// Also show warnings under a fail group
			for _, f := range g.findings {
				if f.Severity == checks.SeverityWarn {
					fmt.Fprintf(w, "      %s %s\n", yellowWarn.Sprint("⚠"), dim.Sprint(f.Message))
				}
			}
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "  Summary: %s, %s, %s\n",
		greenCheck.Sprintf("%d passed", passed),
		yellowWarn.Sprintf("%d warning(s)", warnings),
		redFail.Sprintf("%d failed", failed))
	fmt.Fprintln(w)
}

// worstSeverity returns the most severe finding in the group.
func worstSeverity(findings []checks.Finding) string {
	worst := checks.SeverityPass
	for _, f := range findings {
		if f.Severity == checks.SeverityFail {
			return checks.SeverityFail
		}
		if f.Severity == checks.SeverityWarn {
			worst = checks.SeverityWarn
		}
	}
	return worst
}

// countSeverity counts findings with the given severity.
func countSeverity(findings []checks.Finding, severity string) int {
	n := 0
	for _, f := range findings {
		if f.Severity == severity {
			n++
		}
	}
	return n
}
