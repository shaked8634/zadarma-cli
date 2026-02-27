# Zadarma CLI

A lightweight Go command-line interface for the Zadarma VoIP API.

## Features

- Balance, SIP, DIDs, SMS, PBX, Statistics, Webhooks
- JSON output with `--json` flag
- Debug mode with `-d` / `--debug`
- Minimal dependencies

## Installation

```bash
go build -o zadarma ./cmd/zadarma
```

## Quick Start

### 1. Set up credentials

```bash
export ZADARMA_API_KEY="your_api_key"
export ZADARMA_API_SECRET="your_api_secret"
```

### 2. Check your balance

```bash
./zadarma balance
```

## Common Use Cases

### Send an SMS

**Step 1: Find your phone number**
```bash
./zadarma did list
```
Output:
```
DID: +14155551234 (Type: mobile)
DID: +442071234567 (Type: landline)
```

**Step 2: Send the SMS**
```bash
./zadarma sms send --phone "+14155559999" --message "Hello from Zadarma CLI!"
```

### Check account balance

```bash
./zadarma balance
# Balance: 123.45 USD
```

### List your SIP accounts

```bash
./zadarma sip list
```

### Debug API requests

```bash
./zadarma -d did list
```
Shows:
```
[DEBUG] Request: GET https://api.zadarma.com/v1/did/
[DEBUG] Authorization: 1c64dee7ee76...
[DEBUG] Response: HTTP 200 (523 bytes)
```

## Commands

| Command | Description |
|---------|-------------|
| `balance` | Get account balance |
| `sip list` | List SIP accounts |
| `did list` | List phone numbers (DIDs) |
| `sms send --phone <num> --message <text>` | Send SMS |
| `pbx info` | Get PBX configuration |
| `statistics` | Get call statistics |
| `webhook set <url>` | Set webhook URL |
| `webhook get` | Get current webhook URL |
| `webhook listen` | Start local webhook listener |

## Global Flags

| Flag | Description |
|------|-------------|
| `-k, --key` | API key (or set `ZADARMA_API_KEY`) |
| `-s, --secret` | API secret (or set `ZADARMA_API_SECRET`) |
| `--json` | Output in JSON format |
| `-d, --debug` | Enable debug output |
| `-v, --version` | Show version |

## API Documentation

https://zadarma.com/en/support/api/

## Testing

```bash
go test ./...
```

## License

MIT
