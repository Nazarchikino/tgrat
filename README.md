# TgRAT (Telegram Remote Administration Tool)

[![Go Report Card](https://goreportcard.com/badge/github.com/nazarchikino/tgrat)](https://goreportcard.com/report/github.com/nazarchikino/tgrat)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/nazarchikino/tgrat)](https://golang.org)

**TgRAT** is a lightweight, cross-platform remote administration tool written in Go that uses the Telegram Bot API as a secure, bidirectional Command and Control (C2) channel. It allows remote management, command execution, and host monitoring directly via Telegram chat interfaces.

---

## ⚠️ Disclaimer

> **NOTICE:** This software is designed and distributed strictly for educational purposes, authorized security assessments, and legitimate remote systems administration. Unauthorized access or deployment to systems without explicit, written consent from the owner is strictly prohibited and violates local and international laws. The developer accepts no liability for any misuse or damages caused by this tool.

---

## ✨ Features

- **Telegram Bot C2**: Transmit commands and receive stdout/stderr securely over HTTPS via Telegram's infrastructure.
- **Cross-Platform**: Compiles into standalone binaries for Linux, Windows, and macOS without external runtime dependencies.
- **Remote Command Execution**: Run arbitrary shell/terminal commands and view realtime output directly in chat.
- **Telemetry & Host Recon**: Inspect system properties such as hostname, running user, environment variables, and network configurations.
- **File Transfer**: Send files from the remote host to Telegram or upload files directly into the target filesystem.
- **Stealth & Portability**: Zero dependencies, compact binary size, and configurable background execution.

---

## 📋 Prerequisites

- **Go**: Version `1.21` or newer.
- **Telegram Bot Token**: Generated via [@BotFather](https://t.me/BotFather).
- **Telegram User ID**: Your authorized numerical user ID (can be fetched using [@userinfobot](https://t.me/userinfobot)) to prevent unauthorized access.

---

## 🚀 Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/nazarchikino/tgrat.git
cd tgrat
```

### 2. Configure Credentials

Export your Telegram Bot credentials as environment variables, or insert them into your build configuration:

```bash
export TG_BOT_TOKEN="your_bot_token_here"
export TG_CHAT_ID="your_authorized_chat_id"
```

### 3. Dependencies & Compilation

Download and verify all project dependencies:

```bash
go mod tidy
```

Compile a standard binary:

```bash
go build -ldflags="-s -w" -o tgrat main.go
```

#### Cross-Compilation Commands

- **Windows (x64, Hidden Console):**
  ```bash
  GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o tgrat.exe main.go
  ```

- **Linux (x64):**
  ```bash
  GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o tgrat-linux main.go
  ```

- **macOS (Apple Silicon / ARM64):**
  ```bash
  GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o tgrat-darwin main.go
  ```

---

## 🕹️ Command Reference

Interact with the agent directly in the designated Telegram chat:

| Command | Description |
| :--- | :--- |
| `/start` or `/help` | Display the list of available commands and bot status |
| `/info` | Retrieve machine details (OS, user privileges, public/local IP, architecture) |
| `/exec <cmd>` | Execute a system shell command and return the combined output |
| `/download <path>` | Send a file from the host filesystem as a Telegram document |
| `/upload` | Save a file sent in the Telegram chat to the host |
| `/cd <dir>` | Change current working directory on the host |
| `/exit` | Terminate the remote agent process |

---

## 📂 Project Structure

```text
tgrat/
├── LICENSE         # License terms and conditions
├── README.md       # Project overview and instructions
├── go.mod          # Go module path and dependency tracking
├── go.sum          # Cryptographic checksums of direct/indirect dependencies
└── main.go         # Application entry point, Telegram polling loop, and handlers
```

---

## 🛡️ Security Best Practices

- **Restrict Access**: Always hardcode or validate authorized `chat_id` values to prevent third parties from executing commands if they discover your bot username.
- **Token Hygiene**: Never commit your active Telegram Bot tokens or sensitive credentials into Git version control. Use `.gitignore` for local configuration files.

---

## 📄 License

Distributed under the terms of the [MIT License](LICENSE).
