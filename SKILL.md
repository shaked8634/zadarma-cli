---
name: zadarma-cli
description: OpenClaw skill for the Zadarma VoIP API CLI tool. Use this skill to query Zadarma account state and configure SMS webhooks.
---

# Zadarma CLI Skill

This skill describes the current `zadarma-cli` command set in this repository.

## Binary

Build and run:

```bash
go build -o zadarma-cli ./cmd/zadarma
./zadarma-cli --help
```

## Authentication

Prefer environment variables:

```bash
export ZADARMA_API_KEY="your_api_key"
export ZADARMA_API_SECRET="your_api_secret"
```

Flags also work:

```bash
./zadarma-cli --key "KEY" --secret "SECRET" balance
```

## Agent Output Rule

When using this CLI as an agent, always request machine-readable output:

```bash
./zadarma-cli --output json <command>
```

## Commands

### `balance`

```bash
./zadarma-cli --output json balance
```

### `sip`

```bash
./zadarma-cli --output json sip list
./zadarma-cli --output json sip info 123456
./zadarma-cli --output json sip caller-id --id 123456 --number "+14155551234"
```

### `phone`

```bash
./zadarma-cli --output json phone list
./zadarma-cli --output json phone list 14155551234
./zadarma-cli --output json phone countries
./zadarma-cli --output json phone country US
./zadarma-cli --output json phone number 14155551234
```

### `sms`

```bash
./zadarma-cli --output json sms send --phone "14155559999" --message "Hello" --sender "MyBrand"
./zadarma-cli --output json sms senders --phones "14155559999,442071234567"
./zadarma-cli --output json sms get-webhook
./zadarma-cli --output json sms set-webhook "https://example.com/webhook" --port 8080
./zadarma-cli --output json sms listen --webhook "https://example.com/webhook" --port 8080
```

### `pbx`

```bash
./zadarma-cli --output json pbx info
./zadarma-cli --output json pbx info --pbx-id "123" --numbers "14155551234,442071234567"
```

### `statistics`

```bash
./zadarma-cli --output json statistics
./zadarma-cli --output json statistics --start "2026-03-01 00:00:00" --end "2026-03-07 23:59:59" --sip 123456
```

### `completion`

```bash
./zadarma-cli completion bash
./zadarma-cli completion zsh
./zadarma-cli completion fish
./zadarma-cli completion powershell
```

Pre-generated completion scripts are in `completions/`.

## Capability Summary

- ✅ Account balance
- ✅ SIP list/status
- ✅ SIP caller-id setting
- ✅ Phone/DID listing and lookup
- ✅ SMS send
- ✅ SMS sender lookup
- ✅ SMS webhook get/set
- ✅ SMS local listener
- ✅ PBX info
- ✅ Call statistics
- ✅ Shell completion generation
