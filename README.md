# RepoHealth

A fast Go CLI tool that scans any local repository and reports on its "health" — README quality, broken internal markdown links, oversized committed files, and missing standard files.

Point it at any directory and get an honest report on whether it looks like a well-maintained repo.

## Installation

Choose **one** of the following installation methods:

### Option 1: From Source (Clone and Build)

Use this option if you want to inspect the source code or build it manually:

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

### Option 2: Using `go install` (Fastest)

Use this option to install it instantly without manually cloning the folder:

```bash
go install github.com/DSiddharth24/RepoHealth@latest
```

This installs the binary directly to your `$GOPATH/bin` (or `$HOME/go/bin`). Make sure that folder is in your environment `PATH`.

## Usage

> [!IMPORTANT]
> **Windows PowerShell/CMD Users:** If you built the tool from source and are running it locally in the folder, you must use `.\repohealth` or `.\repohealth.exe` instead of just `repohealth`.
>
> Examples below assume a global installation. If using a local binary, prefix them with `.\`.

```bash
# Scan the current directory
repohealth

# Scan a specific repository folder
repohealth /path/to/your-repository

# JSON output (for CI pipelines or scripting)
repohealth --json /path/to/your-repository

# Strict mode — treat warnings as failures (great for CI gating)
repohealth --strict /path/to/your-repository

# Disable colors (for log files or CI output)
repohealth --no-color /path/to/your-repository

# Ignore additional patterns (repeatable)
repohealth --ignore "*.generated.go" --ignore "testdata" /path/to/your-repository

# Combine flags
repohealth --strict --json --no-color /path/to/your-repository
```

### Quick Start Walkthrough

Here is a step-by-step guide on how to scan a different folder on your machine:

1. **Open your terminal** (Terminal, PowerShell, or Command Prompt).
2. **Execute the tool** and point it to the repository folder you want to inspect:
   ```bash
   # On macOS/Linux:
   repohealth /path/to/your-repository

   # On Windows:
   repohealth C:\path\to\your-repository
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

## How to Fix Health Issues (Your Action Checklist)

If your repository gets warnings or failures, here is exactly what you need to do to resolve them step-by-step:

1. **Fixing README issues:**
   - **If Missing:** Create a `README.md` file in the root folder of your project.
   - **If it's a Stub (less than 150 words):** Expand the README to explain what the project is, how to install it, and how to use it.
   - **If Unstructured (no headings):** Use markdown headers (like `# Project Name` or `## Installation`) to organize your content.

2. **Fixing Broken Links:**
   - Look at the files flagged in the report (e.g., `README.md → ./LICENSE`).
   - Check if you have typos in the relative links or if you referenced files that were deleted or renamed.
   - Correct the link target path relative to the file containing it, or create the missing file.

3. **Fixing Large Files:**
   - **If a file is over 25MB (Failure) or over 5MB (Warning):** Determine if this file should be part of the repository history (e.g. large assets, binaries, database dumps, node modules, build files).
   - **To Resolve:** Add the pattern or file name to your `.gitignore` file (e.g., `*.mp4`, `build/`, `node_modules/`). Once ignored, `RepoHealth` will skip checking it.

4. **Fixing Standard Files:**
   - **If LICENSE is missing:** Create a file named `LICENSE`, `LICENSE.md`, or `LICENSE.txt` at the root and fill it with your repository's license text (e.g., MIT, Apache 2.0).
   - **If .gitignore is missing:** Create a `.gitignore` file at the root to specify which temporary build files, logs, or dependency folders git should ignore.

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
