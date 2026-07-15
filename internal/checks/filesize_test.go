package checks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DSiddharth24/RepoHealth/internal/scanner"
)

func TestFileSizeCheck_Pass(t *testing.T) {
	root := t.TempDir()

	// Create small files (well under 5MB)
	os.WriteFile(filepath.Join(root, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(root, "README.md"), []byte("# Hello"), 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &FileSizeCheck{}
	findings := check.Run(root, sc)

	if len(findings) != 1 {
		t.Fatalf("Expected 1 finding, got %d: %+v", len(findings), findings)
	}
	if findings[0].Severity != SeverityPass {
		t.Errorf("Expected pass, got %s: %s", findings[0].Severity, findings[0].Message)
	}
}

func TestFileSizeCheck_Warn(t *testing.T) {
	root := t.TempDir()

	// Create a file between 5MB and 25MB (6MB)
	size := 6 * 1024 * 1024
	data := make([]byte, size)
	os.WriteFile(filepath.Join(root, "medium.bin"), data, 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &FileSizeCheck{}
	findings := check.Run(root, sc)

	hasWarn := false
	for _, f := range findings {
		if f.Severity == SeverityWarn && strings.Contains(f.Message, "medium.bin") {
			hasWarn = true
		}
	}
	if !hasWarn {
		t.Errorf("Expected warning for 6MB file, got: %+v", findings)
	}
}

func TestFileSizeCheck_Fail(t *testing.T) {
	root := t.TempDir()

	// Create a file over 25MB (26MB)
	size := 26 * 1024 * 1024
	data := make([]byte, size)
	os.WriteFile(filepath.Join(root, "huge.bin"), data, 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &FileSizeCheck{}
	findings := check.Run(root, sc)

	hasFail := false
	for _, f := range findings {
		if f.Severity == SeverityFail && strings.Contains(f.Message, "huge.bin") {
			hasFail = true
		}
	}
	if !hasFail {
		t.Errorf("Expected fail for 26MB file, got: %+v", findings)
	}
}

func TestFileSizeCheck_RespectsGitignore(t *testing.T) {
	root := t.TempDir()

	// Create .gitignore that ignores the large file
	os.WriteFile(filepath.Join(root, ".gitignore"), []byte("*.bin\n"), 0644)

	// Create a 6MB file that should be ignored
	size := 6 * 1024 * 1024
	data := make([]byte, size)
	os.WriteFile(filepath.Join(root, "ignored.bin"), data, 0644)

	sc, _ := scanner.NewScanner(root, nil)
	check := &FileSizeCheck{}
	findings := check.Run(root, sc)

	// Should pass because the large file is gitignored
	if len(findings) != 1 {
		t.Fatalf("Expected 1 finding, got %d: %+v", len(findings), findings)
	}
	if findings[0].Severity != SeverityPass {
		t.Errorf("Expected pass (large file gitignored), got %s: %s", findings[0].Severity, findings[0].Message)
	}
}
