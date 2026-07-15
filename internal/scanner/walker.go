package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// Scanner walks a repository tree while respecting .gitignore patterns
// and user-supplied ignore globs.
type Scanner struct {
	Root          string
	gitIgnore     *ignore.GitIgnore
	extraPatterns []string
	Warnings      []string
}

// NewScanner creates a scanner rooted at the given path.
// It parses .gitignore if present and accepts additional glob patterns to skip.
func NewScanner(root string, extraIgnores []string) (*Scanner, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("cannot access path %q: %w", root, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", root)
	}

	s := &Scanner{
		Root:          root,
		extraPatterns: extraIgnores,
	}

	gitignorePath := filepath.Join(root, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		gi, err := ignore.CompileIgnoreFile(gitignorePath)
		if err != nil {
			s.Warnings = append(s.Warnings,
				fmt.Sprintf("Failed to parse .gitignore: %v — scanning everything", err))
		} else {
			s.gitIgnore = gi
		}
	}

	return s, nil
}

// Walk traverses the repository tree, calling fn for every file that is not
// ignored by .gitignore, extra patterns, or the .git directory.
// Permission errors are logged as warnings and skipped, not fatal.
func (s *Scanner) Walk(fn func(path string, info fs.FileInfo) error) error {
	return filepath.WalkDir(s.Root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Permission denied or other access error — warn and skip
			if os.IsPermission(err) {
				s.Warnings = append(s.Warnings,
					fmt.Sprintf("Permission denied: %s (skipped)", path))
				if d != nil && d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
			return err
		}

		// Compute the relative path for matching
		rel, relErr := filepath.Rel(s.Root, path)
		if relErr != nil {
			rel = path
		}
		// Normalise to forward slashes for consistent matching
		relSlash := filepath.ToSlash(rel)

		// Always skip .git directory
		if d.IsDir() && d.Name() == ".git" {
			return fs.SkipDir
		}

		// Check .gitignore patterns
		if s.gitIgnore != nil && rel != "." {
			matchPath := relSlash
			if d.IsDir() {
				matchPath += "/"
			}
			if s.gitIgnore.MatchesPath(matchPath) {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
		}

		// Check extra ignore patterns
		if s.shouldIgnoreExtra(relSlash, d.IsDir()) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		// Skip directories themselves — only invoke fn for files
		if d.IsDir() {
			return nil
		}

		info, infoErr := d.Info()
		if infoErr != nil {
			if os.IsPermission(infoErr) {
				s.Warnings = append(s.Warnings,
					fmt.Sprintf("Permission denied reading info: %s (skipped)", path))
				return nil
			}
			return infoErr
		}

		return fn(path, info)
	})
}

// shouldIgnoreExtra checks whether the given relative path matches any of
// the user-supplied extra ignore glob patterns.
func (s *Scanner) shouldIgnoreExtra(relSlash string, isDir bool) bool {
	for _, pattern := range s.extraPatterns {
		// Match against the full relative path
		if matched, _ := filepath.Match(pattern, relSlash); matched {
			return true
		}
		// Match against just the base name
		base := filepath.Base(relSlash)
		if matched, _ := filepath.Match(pattern, base); matched {
			return true
		}
		// For directories, also try matching with trailing slash stripped from pattern
		if isDir {
			p := strings.TrimSuffix(pattern, "/")
			if matched, _ := filepath.Match(p, base); matched {
				return true
			}
		}
	}
	return false
}

// CollectFiles returns all non-ignored file paths (absolute) in the repo.
func (s *Scanner) CollectFiles() ([]string, error) {
	var files []string
	err := s.Walk(func(path string, info fs.FileInfo) error {
		files = append(files, path)
		return nil
	})
	return files, err
}

// CollectMDFiles returns absolute paths to all non-ignored .md files.
func (s *Scanner) CollectMDFiles() ([]string, error) {
	var files []string
	err := s.Walk(func(path string, info fs.FileInfo) error {
		if strings.EqualFold(filepath.Ext(path), ".md") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
