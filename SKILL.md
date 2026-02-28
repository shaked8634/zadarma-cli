---
name: zadarma-cli
description: OpenClaw skill for the Zadarma VoIP API CLI tool. Use this skill to interact with Zadarma's VoIP services via command line.
---

# Zadarma CLI Skill

This skill provides OpenClaw agents with knowledge about the `zadarma-cli` tool, a Go-based CLI for the Zadarma VoIP API.

## 📦 Installation

The `zadarma` binary is included in this repository (built for Linux amd64). No separate installation needed.

```bash
# Clone the repository
git clone https://forgejo.o-st.dev/voidclaw/zadarma-cli.git
cd zadarma-cli

# The binary is already built
./zadarma --help
```

## 🔐 Authentication

Set your Zadarma API credentials as environment variables:

```bash
export ZADARMA_API_KEY="your_api_key"
export ZADARMA_API_SECRET="your_api_secret"
```

Or pass them via flags:

```bash
./zadarma --key "KEY" --secret "SECRET" balance
```

## 🚀 Usage

### Basic Commands

**Check balance:**
```bash
./zadarma balance
```

**List phone numbers (DIDs):**
```bash
./zadarma did list
```

**List SIP accounts:**
```bash
./zadarma sip list
```

**Send SMS:**
```bash
./zadarma sms send --phone "+1234567890" --message "Hello from Zadarma CLI"
```

**Get PBX info:**
```bash
./zadarma pbx info
```

### Webhook Management

**Set webhook URL:**
```bash
./zadarma webhook set "https://your-domain.com/webhook"
```

**Listen for SMS webhooks (local testing):**
```bash
./zadarma webhook listen --port 8080
```
*Note: You'll need to expose the port (e.g., via ngrok) and set the public URL with `webhook set`.*

### JSON Output

Add `--json` flag for machine-readable output:

```bash
./zadarma did list --json
```

## 📚 API Coverage

- ✅ Balance inquiry
- ✅ SIP account management
- ✅ DID (phone number) management
- ✅ SMS sending
- ✅ PBX information
- ✅ Webhook configuration & listening
- ✅ Statistics retrieval

## 🛠️ Development

See `NOTES.md` for architecture details and `README.md` for quick start.

## 🐧 Notes

- This skill is embedded in the repository—no separate wrapper script.
- The CLI follows Unix philosophy: simple, composable, pipe-friendly.
- All API calls use Zadarma's HMAC-SHA1 authentication.

---

Built with Go • Open source • Forgejo CI/CD