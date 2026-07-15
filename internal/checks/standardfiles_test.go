package checks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/DSiddharth24/RepoHealth/internal/scanner"
)

func TestStandardFilesCheck_Pass(t *testing.T) {
	root := t.TempDir()

	os.WriteFile(filepath.Join(root, "LICENSE"), []byte("MIT License"), 0644)
	os.WriteFile(filepath.Join(root, ".gitignore"), []byte("*.log\n"), 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &StandardFilesCheck{}
	findings := check.Run(root, sc)

	if len(findings) != 1 {
		t.Fatalf("Expected 1 finding, got %d: %+v", len(findings), findings)
	}
	if findings[0].Severity != SeverityPass {
		t.Errorf("Expected pass, got %s: %s", findings[0].Severity, findings[0].Message)
	}
}

func TestStandardFilesCheck_MissingLicense(t *testing.T) {
	root := t.TempDir()

	os.WriteFile(filepath.Join(root, ".gitignore"), []byte("*.log\n"), 0644)
	// No LICENSE file

	sc, _ := scanner.NewScanner(root, nil)
	check := &StandardFilesCheck{}
	findings := check.Run(root, sc)

	hasLicenseWarn := false
	for _, f := range findings {
		if f.Severity == SeverityWarn && f.Message == "No LICENSE file found" {
			hasLicenseWarn = true
		}
	}
	if !hasLicenseWarn {
		t.Errorf("Expected warning about missing LICENSE, got: %+v", findings)
	}
}

func TestStandardFilesCheck_MissingGitignore(t *testing.T) {
	root := t.TempDir()

	os.WriteFile(filepath.Join(root, "LICENSE"), []byte("MIT License"), 0644)
	// No .gitignore file

	sc, _ := scanner.NewScanner(root, nil)
	check := &StandardFilesCheck{}
	findings := check.Run(root, sc)

	hasGitignoreWarn := false
	for _, f := range findings {
		if f.Severity == SeverityWarn && f.Message == "No .gitignore file found" {
			hasGitignoreWarn = true
		}
	}
	if !hasGitignoreWarn {
		t.Errorf("Expected warning about missing .gitignore, got: %+v", findings)
	}
}

func TestStandardFilesCheck_BothMissing(t *testing.T) {
	root := t.TempDir()
	// No LICENSE, no .gitignore

	sc, _ := scanner.NewScanner(root, nil)
	check := &StandardFilesCheck{}
	findings := check.Run(root, sc)

	warnCount := 0
	for _, f := range findings {
		if f.Severity == SeverityWarn {
			warnCount++
		}
	}
	if warnCount != 2 {
		t.Errorf("Expected 2 warnings (both missing), got %d: %+v", warnCount, findings)
	}
}

func TestStandardFilesCheck_LicenseMD(t *testing.T) {
	root := t.TempDir()

	// LICENSE.md variant should also count
	os.WriteFile(filepath.Join(root, "LICENSE.md"), []byte("# MIT License"), 0644)
	os.WriteFile(filepath.Join(root, ".gitignore"), []byte("*.log\n"), 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &StandardFilesCheck{}
	findings := check.Run(root, sc)

	if len(findings) != 1 {
		t.Fatalf("Expected 1 finding, got %d: %+v", len(findings), findings)
	}
	if findings[0].Severity != SeverityPass {
		t.Errorf("Expected pass for LICENSE.md variant, got %s: %s", findings[0].Severity, findings[0].Message)
	}
}
