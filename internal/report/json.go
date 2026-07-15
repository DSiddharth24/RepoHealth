package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/DSiddharth24/RepoHealth/internal/checks"
)

// JSONReport is the structure marshaled to JSON output.
type JSONReport struct {
	Findings []checks.Finding `json:"findings"`
	Summary  JSONSummary      `json:"summary"`
	ExitCode int              `json:"exit_code"`
}

// JSONSummary holds aggregate counts for the JSON report.
type JSONSummary struct {
	Passed   int `json:"passed"`
	Warnings int `json:"warnings"`
	Failed   int `json:"failed"`
}

// JSON writes a machine-readable JSON report to the given writer.
func JSON(w io.Writer, findings []checks.Finding, exitCode int) error {
	summary := JSONSummary{}

	// Count unique checks by their worst severity
	checkWorst := map[string]string{}
	for _, f := range findings {
		cur, ok := checkWorst[f.Check]
		if !ok || severityRank(f.Severity) > severityRank(cur) {
			checkWorst[f.Check] = f.Severity
		}
	}
	for _, worst := range checkWorst {
		switch worst {
		case checks.SeverityPass:
			summary.Passed++
		case checks.SeverityWarn:
			summary.Warnings++
		case checks.SeverityFail:
			summary.Failed++
		}
	}

	report := JSONReport{
		Findings: findings,
		Summary:  summary,
		ExitCode: exitCode,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON report: %w", err)
	}

	_, err = fmt.Fprintln(w, string(data))
	return err
}

// severityRank returns a numeric rank for sorting severities.
func severityRank(s string) int {
	switch s {
	case checks.SeverityPass:
		return 0
	case checks.SeverityWarn:
		return 1
	case checks.SeverityFail:
		return 2
	default:
		return -1
	}
}
