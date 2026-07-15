# RepoHealth

A fast Go CLI tool to scan a repository's health: README quality, broken internal links, oversized files, and missing standard files (license, .gitignore).

---

## 🚀 Quick Start (Choose One)

### Option A: Install Globally (Fastest)
Installs the tool globally. Run it from anywhere.
```bash
go install github.com/DSiddharth24/RepoHealth@latest
```

### Option B: Build Locally
Manually clone and build the executable.
```bash
git clone https://github.com/DSiddharth24/RepoHealth.git
cd RepoHealth
go build -o repohealth.exe .
```

---

## 💻 Usage

> [!IMPORTANT]
> **PowerShell Users:** If running a locally built binary (Option B), you must prefix the command with `.\` (e.g., `.\repohealth.exe`).

```bash
# Scan the current folder
repohealth

# Scan a specific repository
repohealth /path/to/your-repository
```

### Useful Flags
| Flag | Description |
| :--- | :--- |
| `--json` | Output machine-readable JSON report |
| `--strict` | Treat warnings as failures (exits with code 1) |
| `--no-color` | Disable ANSI terminal colors |
| `--ignore <pattern>` | Skip specific files/folders (e.g., `--ignore "vendor/*"`) |

---

## 🔍 Health Checks

| Check | Passes If... | Fails / Warns If... |
| :--- | :--- | :--- |
| **README** | Present, ≥150 words, has headings | ✗ Missing entirely <br> ⚠ < 150 words or no headings |
| **Broken Links** | All internal markdown links resolve | ✗ Relative target file does not exist |
| **Large Files** | Tracked files are under 5 MB | ⚠ File is 5–25 MB <br> ✗ File is > 25 MB (respects `.gitignore`) |
| **Standard Files** | `LICENSE` and `.gitignore` exist | ⚠ Either file is missing |

---

## 🛠️ Action Checklist (How to Pass)

If your repository gets warnings or failures, fix them using this guide:

*   **README:** Create `README.md` at root. Add structured headers (e.g. `# My Project`) and write at least 150 words.
*   **Broken Links:** Fix the target paths in your markdown links `[text](./link)` so they exist relative to the file.
*   **Large Files:** If you have large files, add them or their folders to your `.gitignore` to skip them.
*   **Standard Files:** Place a `LICENSE` file and a `.gitignore` file in your root directory.

---

## ❓ Troubleshooting

*   **Error:** `destination path 'RepoHealth' already exists...`
    *   **Fix:** You already have a folder named `RepoHealth`. Delete/rename it or clone to a different name:
        ```bash
        git clone https://github.com/DSiddharth24/RepoHealth.git MyRepoHealth
        ```
*   **Error:** `go: cannot find main module...`
    *   **Fix:** You ran `go build` inside the wrong folder. Make sure to `cd RepoHealth` first.
*   **Error:** `The term 'repohealth' is not recognized...`
    *   **Fix:** On Windows PowerShell, type `.\repohealth.exe` instead of `repohealth` when executing local files.

---

## License
MIT
