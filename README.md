# Zadarma CLI

A Go command-line interface for the Zadarma VoIP API.

## Features

- Account balance lookup
- SIP listing and status lookup
- Phone/DID management commands:
  - List owned numbers
  - List available countries
  - List country destinations
  - Inspect a specific virtual number
- SMS commands:
  - Send SMS
  - List valid senders
  - Get webhook URL
  - Set webhook URL and enable SMS hooks
  - Run a local listener for incoming SMS events
- PBX info lookup (with optional filters)
- Call statistics lookup (with optional filters)
- Text (default) and JSON output (`--output text|json`)
- Shell completion support (`completion` command)
- Debug logging (`--debug`)

## Installation

Build from source:

```bash
go build -o zadarma-cli ./cmd/zadarma
```

## Authentication

Use environment variables (recommended):

```bash
export ZADARMA_API_KEY="your_api_key"
export ZADARMA_API_SECRET="your_api_secret"
```

Or pass per command:

```bash
zadarma-cli --key "KEY" --secret "SECRET" balance
```

## Usage

```bash
zadarma-cli [command] [flags]
```

### Global Flags

| Flag | Description |
|------|-------------|
| `-h, --help` | Show help |
| `-v, --version` | Show CLI version |
| `-k, --key` | Zadarma API key |
| `-s, --secret` | Zadarma API secret |
| `-o, --output` | Output format: `text` (default) or `json` |
| `-d, --debug` | Enable debug output |

## Command Reference

### `balance`
- `zadarma-cli balance`

### `sip`
- `zadarma-cli sip list`
- `zadarma-cli sip info <ID>`

### `phone`
- `zadarma-cli phone list [number...]`
- `zadarma-cli phone countries`
- `zadarma-cli phone country <code>`
- `zadarma-cli phone number <number>`

### `sms`
- `zadarma-cli sms send --phone <number> --message <text> [--sender <sender>]`
- `zadarma-cli sms senders [phones]`
- `zadarma-cli sms senders --phones <comma-separated>`
- `zadarma-cli sms get-webhook`
- `zadarma-cli sms set-webhook <WEBHOOK> [--port <port>]`
- `zadarma-cli sms listen [--webhook <WEBHOOK>] [--port <port>]`

### `pbx`
- `zadarma-cli pbx info [--pbx-id <id>] [--numbers <comma-separated>]`

### `statistics`
- `zadarma-cli statistics [--start "YYYY-MM-DD HH:MM:SS"] [--end "YYYY-MM-DD HH:MM:SS"] [--sip <id>]`

### `completion`
- `zadarma-cli completion bash`
- `zadarma-cli completion zsh`
- `zadarma-cli completion fish`
- `zadarma-cli completion powershell`

Pre-generated scripts are available in `completions/`.

## Example: JSON Output

```bash
zadarma-cli --output json phone list
```

## API Documentation

https://zadarma.com/en/support/api/

## License

MIT
