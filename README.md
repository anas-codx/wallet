# Secure Vault CLI 🔒

A production-grade, cryptographically secure CLI tool for managing secrets and passwords. Built with Go, Bubble Tea (TUI), and Lip Gloss.

## 🌟 Features

- **End-to-End Encryption**: Uses **AES-256-GCM** for authenticated encryption.
- **Strong Key Derivation**: Passwords are converted to cryptographic keys using **Argon2id**.
- **Interactive TUI**: A beautiful, responsive Terminal User Interface built with Bubble Tea.
- **Fast CLI Flags**: Quick access for power users (`-l`, `-a`, `-v`, `-d`, `-g`).
- **Secure Password Generator**: Generate high-entropy passwords with `ctrl+g` (TUI) or `-g` (CLI).
- **Clipboard Integration**: Instantly copy secrets for immediate use.
- **Persistent Storage**: Secrets are stored in platform-standard configuration directories (e.g., `~/.config/wallet/`).

## 🚀 Installation

Ensure you have [Go](https://go.dev/) installed.

```bash
# Clone the repository
git clone https://github.com/anas-codx/wallet.git
cd wallet

# Build the binary
go build -o wallet main.go
```

## 🛠️ Usage

### Interactive Mode (TUI)
Simply run the binary to enter the vault:
```bash
./wallet
```

### CLI One-Liners
| Command | Description |
| :--- | :--- |
| `./wallet -l` | List all secret keys |
| `./wallet -a <key>` | Add or update a secret |
| `./wallet -v <key>` | View and copy a secret to clipboard |
| `./wallet -d <key>` | Delete a secret |
| `./wallet -g <length>` | Generate and copy a random password |

## 🏗️ Architecture & Design

The project follows **SOLID** principles and **Clean Architecture** patterns:
- **Repository Pattern**: Encapsulated storage logic in `internal/vault`.
- **Strategy Pattern**: Abstracted platform operations (like clipboard) for testability.
- **Dependency Injection**: Decoupled UI from storage logic.
- **Factory Pattern**: Centralized store creation.

## 🔒 Security

- **No Local History**: All secret inputs are masked in the terminal.
- **Safe Defaults**: All passwords generated use a mix of uppercase, lowercase, digits, and symbols.
- **Argon2 Parameters**: Configured for high resistance against GPU/ASIC brute-force attacks.

## 🧪 Testing

Run the full test suite:
```bash
go test -v ./...
```

## 📝 License
MIT
