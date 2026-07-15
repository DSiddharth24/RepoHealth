package checks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/DSiddharth24/RepoHealth/internal/scanner"
)

func TestLinksCheck_Pass(t *testing.T) {
	root := t.TempDir()

	// Create target files
	docsDir := filepath.Join(root, "docs")
	os.MkdirAll(docsDir, 0755)
	os.WriteFile(filepath.Join(docsDir, "setup.md"), []byte("# Setup"), 0644)

	// Create README with valid internal link
	readme := `# Project
See [setup guide](./docs/setup.md) for details.
`
	os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &LinksCheck{}
	findings := check.Run(root, sc)

	if len(findings) != 1 {
		t.Fatalf("Expected 1 finding, got %d: %+v", len(findings), findings)
	}
	if findings[0].Severity != SeverityPass {
		t.Errorf("Expected pass, got %s: %s", findings[0].Severity, findings[0].Message)
	}
}

func TestLinksCheck_BrokenLink(t *testing.T) {
	root := t.TempDir()

	// Create README with a broken internal link (target doesn't exist)
	readme := `# Project
See [missing guide](./docs/nonexistent.md) for details.
`
	os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &LinksCheck{}
	findings := check.Run(root, sc)

	hasFail := false
	for _, f := range findings {
		if f.Severity == SeverityFail {
			hasFail = true
		}
	}
	if !hasFail {
		t.Errorf("Expected fail for broken link, got: %+v", findings)
	}
}

func TestLinksCheck_ExternalLinksIgnored(t *testing.T) {
	root := t.TempDir()

	// Create README with only external links
	readme := `# Project
See [docs](https://example.com) and [more](http://example.org).
`
	os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &LinksCheck{}
	findings := check.Run(root, sc)

	if len(findings) != 1 {
		t.Fatalf("Expected 1 finding, got %d: %+v", len(findings), findings)
	}
	if findings[0].Severity != SeverityPass {
		t.Errorf("Expected pass (externals ignored), got %s: %s", findings[0].Severity, findings[0].Message)
	}
}

func TestLinksCheck_AnchorLinksIgnored(t *testing.T) {
	root := t.TempDir()

	// Create README with only anchor links
	readme := `# Project
See [section](#overview) for details.
`
	os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &LinksCheck{}
	findings := check.Run(root, sc)

	if len(findings) != 1 {
		t.Fatalf("Expected 1 finding, got %d: %+v", len(findings), findings)
	}
	if findings[0].Severity != SeverityPass {
		t.Errorf("Expected pass (anchors ignored), got %s: %s", findings[0].Severity, findings[0].Message)
	}
}

func TestLinksCheck_FragmentStripping(t *testing.T) {
	root := t.TempDir()

	// Create target file
	docsDir := filepath.Join(root, "docs")
	os.MkdirAll(docsDir, 0755)
	os.WriteFile(filepath.Join(docsDir, "guide.md"), []byte("# Guide\n## Section"), 0644)

	// Link with fragment that points to existing file
	readme := `# Project
See [guide section](./docs/guide.md#section) for details.
`
	os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &LinksCheck{}
	findings := check.Run(root, sc)

	if len(findings) != 1 {
		t.Fatalf("Expected 1 finding, got %d: %+v", len(findings), findings)
	}
	if findings[0].Severity != SeverityPass {
		t.Errorf("Expected pass (fragment stripped, file exists), got %s: %s", findings[0].Severity, findings[0].Message)
	}
}

func TestLinksCheck_NoMDFiles(t *testing.T) {
	root := t.TempDir()

	// No markdown files at all
	os.WriteFile(filepath.Join(root, "main.go"), []byte("package main"), 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &LinksCheck{}
	findings := check.Run(root, sc)

	if len(findings) != 1 {
		t.Fatalf("Expected 1 finding, got %d: %+v", len(findings), findings)
	}
	if findings[0].Severity != SeverityPass {
		t.Errorf("Expected pass, got %s: %s", findings[0].Severity, findings[0].Message)
	}
}

func TestLinksCheck_MultipleBrokenLinks(t *testing.T) {
	root := t.TempDir()

	readme := `# Project
See [a](./missing1.md) and [b](./missing2.md).
`
	os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &LinksCheck{}
	findings := check.Run(root, sc)

	failCount := 0
	for _, f := range findings {
		if f.Severity == SeverityFail {
			failCount++
		}
	}
	if failCount != 2 {
		t.Errorf("Expected 2 broken link failures, got %d: %+v", failCount, findings)
	}
}
