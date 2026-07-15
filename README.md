# RepoHealth

A fast Go CLI tool to scan a repository's health: README quality, broken internal links, oversized files, and missing standard files (license, .gitignore).

---

## 📋 Prerequisites

Before installing, you must have **Go (GoLang) 1.22 or higher** installed on your system.

*   Download it from the official site: **[go.dev/doc/install](https://go.dev/doc/install)**
*   Verify it is installed by running this command in your terminal:
    ```bash
    go version
    ```

---

## 🚀 Quick Start (Choose One)

### Option A: Install Globally (Fastest)

Follow these **two mandatory steps** to install globally:

#### Step 1: Install the package
```bash
go install github.com/DSiddharth24/RepoHealth@latest
```

#### Step 2: Add Go Binary directory to PATH
To ensure your shell recognizes the `repohealth` command, add the installation folder to your environment PATH:

*   **Windows (PowerShell):**
    ```powershell
    $env:Path += ";$env:GOPATH\bin"
    ```
*   **Linux / macOS:**
    ```bash
    export PATH=$PATH:$HOME/go/bin
    ```

### Option B: Build Locally
Manually clone and build the executable.

**Linux / macOS:**
```bash
git clone https://github.com/DSiddharth24/RepoHealth.git
cd RepoHealth
go build -o repohealth .
```

**Windows:**
```powershell
git clone https://github.com/DSiddharth24/RepoHealth.git
cd RepoHealth
go build -o repohealth.exe .
```

---

## 💻 Usage

### 1. If Installed Globally (Option A)
Run the tool globally from any folder without prefixes:
```bash
# Scan the current folder
repohealth

# Scan a specific repository
repohealth /path/to/your-repository
#Example for PATH
repohealth /Users/siddh/Gemini-Chatbot
```

### 2. If Built Locally (Option B)
You must specify the path to the local executable file in the current folder:

*   **Linux / macOS:**
    ```bash
    ./repohealth /Users/siddh/Gemini-Chatbot
    ```
*   **Windows (PowerShell):**
    ```powershell
    .\repohealth.exe C:\Users\siddh\Gemini-Chatbot
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
