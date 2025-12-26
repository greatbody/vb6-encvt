# VB6 Encoding Converter (vb6enc)

![Vibe Coded](https://img.shields.io/badge/Vibe-Coded-ff69b4)

A lightweight, robust command-line tool designed to bridge the gap between legacy VB6 environments (GBK/GB2312) and modern AI/LLM workflows (UTF-8).

## ⚡️ Vibe Coded
This project was built with a "Vibe Coded" philosophy: simple, effective, and gets the job done without over-engineering. It solves exactly one problem: making old code talk to new AI.

## 🚀 Why?

Modern AI coding assistants and Agentic workflows thrive in a UTF-8 world. Legacy VB6 projects, however, are deeply rooted in GBK (Simplified Chinese) encoding. 

This tool allows you to:
1.  **Batch Convert** your entire source tree to UTF-8 for AI processing / LLM refactoring.
2.  **Batch Restore** the files back to GBK so the VB6 IDE allows you to compile and edit without "" errors.

## ✨ Features

- **🛡 Safe**: Performs atomic file replacements using temporary files. No half-written corruptions.
- **🧠 Smart**: Uses heuristic detection to distinguish text files from binaries. Automatically ignores `.git`, `.svn`, `bin`, `obj`.
- **⚡️ Fast**: Written in pure Go. zero-dependency runtime (single binary).
- **🌍 Cross-Platform**: Compatible with Windows (including Windows 7), macOS, and Linux.

## 🛠 Installation

### Build from Source

```bash
# Clone the repo
git clone https://github.com/greatbody/vb6-encvt.git
cd vb6-encvt

# Build
go build -o vb6enc main.go

# (Optional) Cross-compile for Windows from Mac/Linux
GOOS=windows GOARCH=amd64 go build -o vb6enc.exe main.go
```

## 📖 Usage

### 1. Scan Project Encodings
Check the current state of your project files. This purely transparent and modifies nothing.

```bash
./vb6enc scan /path/to/vb6-project
```

### 2. Convert to UTF-8
Ready for AI ingestion? Convert everything to UTF-8.

```bash
./vb6enc to-utf8 /path/to/vb6-project
```
> *Note: Only files detected as GBK will be converted. Existing UTF-8 files are touched gently and left alone.*

### 3. Convert to GBK
Time to compile in VB6? Convert everything back to GBK.

```bash
./vb6enc to-gb /path/to/vb6-project
```

### 4. Verify Integrity
Ensure all text files are in a known, valid encoding (either UTF-8 or GBK). Useful for spotting corrupted files.

```bash
./vb6enc verify /path/to/vb6-project
```

## 📝 License

MIT
