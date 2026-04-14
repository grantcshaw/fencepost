# fencepost

A CLI tool for managing and rotating API keys across multiple services with audit logging.

---

## Installation

```bash
go install github.com/yourusername/fencepost@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/fencepost.git && cd fencepost && go build -o fencepost .
```

---

## Usage

```bash
# Add a new API key for a service
fencepost add --service stripe --key sk_live_abc123

# List all managed keys
fencepost list

# Rotate a key for a specific service
fencepost rotate --service stripe

# View the audit log
fencepost audit --tail 50
```

Fencepost stores keys securely and maintains a full audit trail of all create, rotate, and delete operations. Configuration is read from `~/.fencepost/config.yaml` by default.

```yaml
# ~/.fencepost/config.yaml
log_path: ~/.fencepost/audit.log
encryption: aes-256
services:
  - stripe
  - github
  - sendgrid
```

---

## Features

- Manage API keys for multiple services from a single CLI
- Automatic key rotation with configurable schedules
- Tamper-evident audit logging for all key operations
- AES-256 encryption for keys at rest

---

## License

MIT © [yourusername](https://github.com/yourusername)