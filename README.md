# RepoHealth

A fast Go CLI tool that scans any local repository and reports on its "health" — README quality, broken internal markdown links, oversized committed files, and missing standard files.

Point it at any directory and get an honest report on whether it looks like a well-maintained repo.

## Installation

### From Source (requires Go 1.22+)

```bash
# Clone the repo
git clone https://github.com/DSiddharth24/RepoHealth.git
cd RepoHealth

# Build the binary
go build -o repohealth .

# (Optional) Move to a directory in your PATH
# Linux/macOS:
sudo mv repohealth /usr/local/bin/
# Windows: move repohealth.exe to a folder in your PATH
```

### Using `go install`

```bash
go install github.com/DSiddharth24/RepoHealth@latest
```

This places the `RepoHealth` binary in your `$GOPATH/bin` (or `$HOME/go/bin`). Make sure that directory is in your `PATH`.

## Usage

```bash
# Scan the current directory
repohealth

# Scan a specific repository folder
repohealth /path/to/other-project

# JSON output (for CI pipelines or scripting)
repohealth --json /path/to/other-project

# Strict mode — treat warnings as failures (great for CI gating)
repohealth --strict /path/to/other-project

# Disable colors (for log files or CI output)
repohealth --no-color /path/to/other-project

# Ignore additional patterns (repeatable)
repohealth --ignore "*.generated.go" --ignore "testdata" /path/to/other-project

# Combine flags
repohealth --strict --json --no-color /path/to/other-project
```

### Quick Start Walkthrough

Here is a step-by-step guide on how to scan a different folder on your machine:

1. **Open your terminal** (Terminal, PowerShell, or Command Prompt).
2. **Execute the tool** and point it to the repository folder you want to inspect:
   ```bash
   # On macOS/Linux:
   repohealth /Users/username/projects/another-repo

   # On Windows:
   repohealth C:\Users\username\projects\another-repo
   ```
3. **Understand the output:**
   - The scanner will recursively walk the target folder.
   - It honors the `.gitignore` inside that project so it does not check dependencies (like `node_modules` or build directories).
   - It will print a report showing whether the target repository passes, warns, or fails on each check.

## What It Checks

| Check | Pass | Warn | Fail |
|-------|------|------|------|
| **README** | Present, ≥150 words, has headings | Exists but < 150 words (stub) or no headings | Missing entirely |
| **Broken Links** | All internal `.md` links resolve | — | Any `[text](./path)` pointing to a nonexistent file |
| **Large Files** | All files under 5 MB | Any file 5–25 MB | Any file over 25 MB |
| **Standard Files** | LICENSE and .gitignore present | Either one missing | — |

## Example Output

### Terminal (default)

```
RepoHealth — scanning ./my-project

  ✓ README                  Present, well-structured (312 words, 6 headings)
  ✗ Broken links            2 issue(s) found
      src/App.tsx → ./docs/architecture.md (target does not exist)
      README.md → ./LICENSE (target does not exist)
  ⚠ Large files              1 warning(s)
      assets/demo.mp4 (18.4 MB)
  ✓ Standard files          LICENSE and .gitignore present

  Summary: 2 passed, 1 warning(s), 1 failed
```

### JSON (`--json`)

```json
{
  "findings": [
    {
      "check": "README",
      "severity": "pass",
      "message": "Present, well-structured (312 words, 6 headings)",
      "file": "README.md"
    },
    {
      "check": "Large files",
      "severity": "warn",
      "message": "assets/demo.mp4 (18.4 MB)",
      "file": "assets/demo.mp4"
    }
  ],
  "summary": { "passed": 3, "warnings": 1, "failed": 0 },
  "exit_code": 0
}
```

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | All checks passed (warnings allowed unless `--strict`) |
| `1` | One or more checks failed, or warnings present with `--strict` |
| `2` | Tool error (invalid path, permission issue) |

## Using in CI (GitHub Actions Example)

```yaml
name: Repo Health Check
on: [push, pull_request]

jobs:
  repohealth:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - run: go install github.com/DSiddharth24/RepoHealth@latest
      - run: RepoHealth --strict --no-color .
```

## Running Tests

```bash
go test ./... -v -count=1
```

## License

MIT — see [LICENSE](./LICENSE) for details.
