# Zadarma CLI

A lightweight, modular Go command-line interface for the Zadarma VoIP API.

## Features

- **Multi-endpoint Support**: Balance, SIP, DIDs, SMS, PBX
- **Cobra CLI Framework**: Intuitive subcommand structure
- **JSON Output**: Use `--json` flag for machine-readable output
- **Unix-Friendly**: Designed for piping and scripting
- **Minimal Dependencies**: Cobra only (proven, stable dependency)
- **Secure**: API credentials via environment variables or flags
- **Tested**: Comprehensive unit tests for all endpoints

## Installation

```bash
go build -o zadarma ./cmd/zadarma
```

## Quick Start

### Setup Credentials

```bash
export ZADARMA_API_KEY="your_api_key"
export ZADARMA_API_SECRET="your_api_secret"
```

### Get Account Balance

```bash
./zadarma balance
```

Output:
```
Balance: 123.45 USD
```

With JSON output:
```bash
./zadarma balance --json
```

## Commands

### Account Info
- **`balance`** — Get account balance
  ```bash
  zadarma balance [--json]
  ```

### SIP Management
- **`sip list`** — List SIP accounts
  ```bash
  zadarma sip list [--json]
  ```

### Phone Numbers (DIDs)
- **`did list`** — List phone numbers
  ```bash
  zadarma did list [--json]
  ```

### SMS
- **`sms send`** — Send an SMS message
  ```bash
  zadarma sms send --phone "<number>" --message "<text>" [--json]
  ```

### PBX
- **`pbx info`** — Get PBX configuration info
  ```bash
  zadarma pbx info [--json]
  ```

## Global Flags

- `--key <key>`: API key (overrides `ZADARMA_API_KEY` env var)
- `--secret <secret>`: API secret (overrides `ZADARMA_API_SECRET` env var)
- `--json`: Output in JSON format
- `-v, --version`: Show version
- `-h, --help`: Show help

## Examples

### Get balance with explicit credentials
```bash
./zadarma --key "abc123" --secret "xyz789" balance
```

### List SIP accounts in JSON
```bash
./zadarma sip list --json
```

### Send SMS
```bash
./zadarma sms send --phone "+14155555555" --message "Hello World"
```

### Pipe balance to grep
```bash
./zadarma balance | grep USD
```

## Architecture

```
cmd/zadarma/          - CLI entry point (cobra root command)
internal/auth/        - HMAC-SHA1 signature generation
internal/client/      - API client (HTTP wrapper with auth)
internal/commands/    - Command handlers for each endpoint
tests/                - Integration tests
```

## API Documentation

See [Zadarma API Docs](https://zadarma.com/en/support/api/)

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./internal/auth -v
go test ./internal/client -v
```

## Contributing

1. Clone the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `go test ./...`
5. Run formatter: `go fmt ./...`
6. Run linter: `go vet ./...`
7. Submit a pull request

## License

MIT
