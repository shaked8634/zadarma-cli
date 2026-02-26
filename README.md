# Zadarma CLI

A lightweight Go command-line interface for the Zadarma VoIP API.

## Features

- **Account Management**: Check account balance
- **Unix-Friendly**: Designed for piping and scripting
- **No Dependencies**: Uses only Go standard library
- **Secure**: API credentials via environment variables or flags

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

## Commands

### `balance`
Retrieves your current account balance.

```bash
zadarma balance
```

### `help`
Shows help information.

```bash
zadarma help
zadarma -h
```

## Flags

- `-key`: Zadarma API key (overrides `ZADARMA_API_KEY` env var)
- `-secret`: Zadarma API secret (overrides `ZADARMA_API_SECRET` env var)
- `-v`: Show version
- `-h`: Show help

## Examples

### Get balance with explicit credentials
```bash
./zadarma -key "abc123" -secret "xyz789" balance
```

### Pipeline usage (future feature)
```bash
./zadarma balance | grep USD
```

## API Documentation

See [Zadarma API Docs](https://zadarma.com/en/support/api/)

## Contributing

1. Clone the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `go test ./...`
5. Submit a pull request

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./internal/auth -v
```

## License

MIT
