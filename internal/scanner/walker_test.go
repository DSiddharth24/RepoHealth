package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestWalker_SkipsGitDir(t *testing.T) {
	root := t.TempDir()

	// Create .git directory with a file inside
	gitDir := filepath.Join(root, ".git")
	os.MkdirAll(gitDir, 0755)
	os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main"), 0644)

	// Create a normal file
	os.WriteFile(filepath.Join(root, "main.go"), []byte("package main"), 0644)

	sc, err := NewScanner(root, nil)
	if err != nil {
		t.Fatalf("NewScanner failed: %v", err)
	}

	files, err := sc.CollectFiles()
	if err != nil {
		t.Fatalf("CollectFiles failed: %v", err)
	}

	for _, f := range files {
		rel, _ := filepath.Rel(root, f)
		if filepath.Base(filepath.Dir(f)) == ".git" || rel == ".git" {
			t.Errorf("Walker should skip .git directory, but found: %s", rel)
		}
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d: %v", len(files), files)
	}
}

func TestWalker_RespectsGitignore(t *testing.T) {
	root := t.TempDir()

	// Create .gitignore
	os.WriteFile(filepath.Join(root, ".gitignore"), []byte("build/\n*.log\n"), 0644)

	// Create files that should be ignored
	buildDir := filepath.Join(root, "build")
	os.MkdirAll(buildDir, 0755)
	os.WriteFile(filepath.Join(buildDir, "output.js"), []byte("compiled"), 0644)
	os.WriteFile(filepath.Join(root, "debug.log"), []byte("log data"), 0644)

	// Create files that should NOT be ignored
	os.WriteFile(filepath.Join(root, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(root, "README.md"), []byte("# Hello"), 0644)

	sc, err := NewScanner(root, nil)
	if err != nil {
		t.Fatalf("NewScanner failed: %v", err)
	}

	files, err := sc.CollectFiles()
	if err != nil {
		t.Fatalf("CollectFiles failed: %v", err)
	}

	relFiles := make([]string, 0, len(files))
	for _, f := range files {
		rel, _ := filepath.Rel(root, f)
		relFiles = append(relFiles, filepath.ToSlash(rel))
	}
	sort.Strings(relFiles)

	// Should include .gitignore, main.go, README.md but not build/ or .log
	expected := []string{".gitignore", "README.md", "main.go"}
	sort.Strings(expected)

	if len(relFiles) != len(expected) {
		t.Fatalf("Expected %d files %v, got %d files %v", len(expected), expected, len(relFiles), relFiles)
	}
	for i, exp := range expected {
		if relFiles[i] != exp {
			t.Errorf("Expected file %q at index %d, got %q", exp, i, relFiles[i])
		}
	}
}

func TestWalker_ExtraIgnorePatterns(t *testing.T) {
	root := t.TempDir()

	os.WriteFile(filepath.Join(root, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(root, "test.txt"), []byte("test"), 0644)
	vendorDir := filepath.Join(root, "vendor")
	os.MkdirAll(vendorDir, 0755)
	os.WriteFile(filepath.Join(vendorDir, "dep.go"), []byte("package dep"), 0644)

	sc, err := NewScanner(root, []string{"*.txt", "vendor"})
	if err != nil {
		t.Fatalf("NewScanner failed: %v", err)
	}

	files, err := sc.CollectFiles()
	if err != nil {
		t.Fatalf("CollectFiles failed: %v", err)
	}

	relFiles := make([]string, 0, len(files))
	for _, f := range files {
		rel, _ := filepath.Rel(root, f)
		relFiles = append(relFiles, filepath.ToSlash(rel))
	}

	if len(relFiles) != 1 || relFiles[0] != "main.go" {
		t.Errorf("Expected only [main.go], got %v", relFiles)
	}
}

func TestWalker_CollectMDFiles(t *testing.T) {
	root := t.TempDir()

	os.WriteFile(filepath.Join(root, "README.md"), []byte("# README"), 0644)
	os.WriteFile(filepath.Join(root, "main.go"), []byte("package main"), 0644)
	docsDir := filepath.Join(root, "docs")
	os.MkdirAll(docsDir, 0755)
	os.WriteFile(filepath.Join(docsDir, "guide.md"), []byte("# Guide"), 0644)
	os.WriteFile(filepath.Join(docsDir, "notes.txt"), []byte("notes"), 0644)

	sc, err := NewScanner(root, nil)
	if err != nil {
		t.Fatalf("NewScanner failed: %v", err)
	}

	mdFiles, err := sc.CollectMDFiles()
	if err != nil {
		t.Fatalf("CollectMDFiles failed: %v", err)
	}

	if len(mdFiles) != 2 {
		relFiles := make([]string, len(mdFiles))
		for i, f := range mdFiles {
			rel, _ := filepath.Rel(root, f)
			relFiles[i] = rel
		}
		t.Errorf("Expected 2 .md files, got %d: %v", len(mdFiles), relFiles)
	}
}

func TestWalker_InvalidPath(t *testing.T) {
	_, err := NewScanner("/nonexistent/path/that/does/not/exist", nil)
	if err == nil {
		t.Error("Expected error for nonexistent path, got nil")
	}
}

func TestWalker_NotADirectory(t *testing.T) {
	root := t.TempDir()
	filePath := filepath.Join(root, "file.txt")
	os.WriteFile(filePath, []byte("hello"), 0644)

	_, err := NewScanner(filePath, nil)
	if err == nil {
		t.Error("Expected error for non-directory path, got nil")
	}
}
