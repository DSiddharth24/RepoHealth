package checks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DSiddharth24/RepoHealth/internal/scanner"
)

func TestReadmeCheck_Pass(t *testing.T) {
	root := t.TempDir()

	// Create a well-structured README with >150 words and headings
	var sb strings.Builder
	sb.WriteString("# My Project\n\n")
	sb.WriteString("## Overview\n\n")
	sb.WriteString("This is a comprehensive project that does many wonderful things. ")
	// Generate enough words
	for i := 0; i < 40; i++ {
		sb.WriteString("Lorem ipsum dolor sit amet consectetur adipiscing elit. ")
	}
	sb.WriteString("\n\n## Installation\n\nRun the install command to get started.\n")

	os.WriteFile(filepath.Join(root, "README.md"), []byte(sb.String()), 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &ReadmeCheck{}
	findings := check.Run(root, sc)

	if len(findings) != 1 {
		t.Fatalf("Expected 1 finding, got %d: %+v", len(findings), findings)
	}
	if findings[0].Severity != SeverityPass {
		t.Errorf("Expected pass, got %s: %s", findings[0].Severity, findings[0].Message)
	}
}

func TestReadmeCheck_Missing(t *testing.T) {
	root := t.TempDir()
	// No README.md created

	sc, _ := scanner.NewScanner(root, nil)
	check := &ReadmeCheck{}
	findings := check.Run(root, sc)

	if len(findings) != 1 {
		t.Fatalf("Expected 1 finding, got %d: %+v", len(findings), findings)
	}
	if findings[0].Severity != SeverityFail {
		t.Errorf("Expected fail, got %s: %s", findings[0].Severity, findings[0].Message)
	}
}

func TestReadmeCheck_WarnStub(t *testing.T) {
	root := t.TempDir()

	// Short README with headings but under 150 words
	content := "# My Project\n\nThis is a short README with very few words.\n"
	os.WriteFile(filepath.Join(root, "README.md"), []byte(content), 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &ReadmeCheck{}
	findings := check.Run(root, sc)

	hasWarn := false
	for _, f := range findings {
		if f.Severity == SeverityWarn && strings.Contains(f.Message, "stub") {
			hasWarn = true
		}
	}
	if !hasWarn {
		t.Errorf("Expected warning about stub README, got: %+v", findings)
	}
}

func TestReadmeCheck_WarnNoHeadings(t *testing.T) {
	root := t.TempDir()

	// Long README but no markdown headings
	var sb strings.Builder
	for i := 0; i < 40; i++ {
		sb.WriteString("Lorem ipsum dolor sit amet consectetur adipiscing elit. ")
	}

	os.WriteFile(filepath.Join(root, "README.md"), []byte(sb.String()), 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &ReadmeCheck{}
	findings := check.Run(root, sc)

	hasWarn := false
	for _, f := range findings {
		if f.Severity == SeverityWarn && strings.Contains(f.Message, "no headings") {
			hasWarn = true
		}
	}
	if !hasWarn {
		t.Errorf("Expected warning about no headings, got: %+v", findings)
	}
}
