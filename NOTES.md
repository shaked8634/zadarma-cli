# Zadarma CLI - Design Decisions & Architecture

## Overview
A Go CLI tool for interacting with the Zadarma VoIP API. Designed to follow Unix philosophy: simple, composable, pipe-friendly.

## Architecture

### Package Structure
```
cmd/zadarma/         - CLI entry point, command handlers
internal/auth/       - HMAC-SHA1 signature generation (Zadarma auth)
internal/client/     - HTTP client wrapper for API calls
tests/               - Integration tests
```

### Key Components

#### 1. **Auth Signer** (`internal/auth/signer.go`)
Implements the Zadarma HMAC-SHA1 signature algorithm:
1. Sort params alphabetically
2. Build query string: `param1=value1&param2=value2`
3. Concatenate: `method + paramsStr + md5(paramsStr)`
4. HMAC-SHA1 with secret key
5. Base64 encode
6. Authorization header: `key:signature`

**Design choice:** Separate package keeps auth logic testable and reusable.

#### 2. **API Client** (`internal/client/client.go`)
Wraps HTTP requests with automatic signature generation:
- Handles JSON unmarshaling
- Provides typed method wrappers (e.g., `GetBalance()`)
- Extensible for new endpoints

**Design choice:** Generic `Get()`, `Post()` methods allow easy addition of new features without code duplication.

#### 3. **CLI Entry Point** (`cmd/zadarma/main.go`)
Command-based interface using Go's `flag` package:
- `balance` - Get account balance
- `help` - Show help
- Supports environment variables: `ZADARMA_API_KEY`, `ZADARMA_API_SECRET`
- Can be extended with subcommands

**Design choice:** Lightweight CLI framework using stdlib `flag` to minimize dependencies. Can migrate to `cobra` if needed later.

## Output Format

- Text output is the default for all commands.
- JSON output is enabled only when explicitly requested via either flag:
  - `--json`
  - `--output=json`

Example:
```
zadarma-cli balance --json
```

## Authentication

### Setup
Export credentials as environment variables:
```bash
export ZADARMA_API_KEY="your_api_key"
export ZADARMA_API_SECRET="your_api_secret"
```

Or pass via CLI flags:
```bash
zadarma-cli -key "..." -secret "..." balance
```

### Authentication & Signing Notes

Zadarma uses HMAC-SHA1 signatures. Our `internal/auth/signer.go` implements the official algorithm and matches the
behavior of the reference TypeScript client.

Key points we discovered and validated against production:

- Always sign using the full versioned path, including leading and trailing slashes. Example: `"/v1/sms/send/"`.
- Canonical parameter string must be built like `application/x-www-form-urlencoded`:
  - Keys sorted alphabetically.
  - Spaces encoded as `+` (NOT `%20`). We rely on `url.Values.Encode()` for this behavior.
- Signature string: `method + paramsStr + md5(paramsStr)`.
- HMAC-SHA1 over the string with the API secret; hex-encode the HMAC result, then Base64-encode that hex string.
- Authorization header format: `API_KEY:BASE64_HEX_HMAC`.

HTTP request conventions per API docs:

- For `POST` and `PUT`, set `Content-Type: application/x-www-form-urlencoded` (or `multipart/form-data`, which we don’t
  currently use).
- For `GET`, parameters are appended to the query string. For non-GET, parameters go in the request body in form-encoded
  format. Our client follows this.

Important: Do NOT force `format=json` globally. We removed the implicit `format=json` parameter. If an endpoint supports
alternate formats and JSON is required for CLI presentation, the command layer decides and adds it explicitly as needed.

## CLI Framework Migration

### v0.2.0 - Cobra Framework
Migrated from stdlib `flag` to `spf13/cobra` for:
- Better subcommand structure
- Consistent flag handling across commands
- Easier future expansion
- Auto-generated shell completions

New command structure:
```
zadarma-cli balance [--json]
zadarma-cli sip list [--json]
zadarma-cli did list [--json]
zadarma-cli sms send --phone <num> --message <msg> [--json]
zadarma-cli pbx info [--json]
```

## CLI Command Naming Convention

To keep the CLI consistent and predictable, we follow a clear pattern:

- Use `list` for commands that return multiple entities (collections).
- Use `info` for commands that return detailed information about a single entity.

This convention is applied across the CLI. Examples:

```
zadarma-cli balance [--json]
zadarma-cli sip list [--json]        # list all SIP accounts
zadarma-cli sip info <ID> [--json]   # get detailed info for a single SIP account
zadarma-cli phone list [--json]      # list DID/phone numbers
zadarma-cli phone countries list     # list available country codes and ISO
zadarma-cli phone country info <CC>  # get number types for country (CC is ISO 3166-1 alpha-2)
zadarma-cli pbx info [--json]
```

Make sure commands are discoverable via `--help` and shell completions.

## Webhook Daemon Mode (SMS)

### Overview

The CLI supports running as a daemon to listen for incoming SMS webhooks from Zadarma. This allows real-time processing
of incoming messages without polling.

### Usage

Start the SMS webhook daemon:

```bash
# Text mode (human-readable table format)
./zadarma-cli sms listen --port 8080

# JSON mode (structured output)
./zadarma-cli --output json sms listen --port 8080

# Background the daemon to use other CLI commands simultaneously
./zadarma-cli sms listen --port 8080 &
./zadarma-cli sms send --phone +1234567890 --message "Hello"
```

### How It Works

1. **Fetch Current Webhook URL**: On startup, `sms listen` calls the API to fetch and display the currently configured
   webhook URL.
2. **Start Local HTTP Server**: Listens on the specified port (default: 8080) for incoming POST requests.
3. **Handle Zadarma Verification**: Responds to Zadarma's `zd_echo` verification queries.
4. **Process SMS Events**:
    - In **text mode**: Prints incoming SMS in a formatted table (FROM, TO, TEXT, TIME).
    - In **JSON mode**: Prints the raw JSON event data.
5. **Run as Daemon**: The daemon runs indefinitely until stopped with Ctrl+C. Can be backgrounded with `&` to keep the
   CLI available for other commands.

### Setup

Before using `sms listen`, you must:

1. **Set the webhook URL** via the API:
   ```bash
   ./zadarma-cli webhook set <WEBHOOK_URL>
   ```

2. **Verify the webhook** in the Zadarma account dashboard (optional).

3. **Start the daemon**:
   ```bash
   ./zadarma-cli sms listen --port 8080
   ```

### Background Execution

To run the daemon in the background while using other CLI commands:

```bash
./zadarma-cli sms listen --port 8080 &
echo $!  # Save the daemon PID

# Now you can use other commands
./zadarma-cli sms send --phone +1234567890 --message "Hi there"

# When done, stop the daemon
kill %1  # or kill <PID>
```

### Output Examples

**Text Mode (default)**:

```
Current webhook URL: https://abc123.ngrok.io
Listening for SMS webhooks on port 8080...
Press Ctrl+C to stop.

--- INCOMING SMS ---
FROM   +1234567890
TO     +1987654321
TEXT   Hello, this is a test
TIME   1709251234
-------------------
```

**JSON Mode** (`--output json`):

```json
{
  "event": "SMS",
  "caller_id": "+1234567890",
  "caller_did": "+1987654321",
  "text": "Hello, this is a test",
  "timestamp": 1709251234
}
```

## Future Enhancements

### Phase 1 (Complete)
- [x] Get balance
- [x] SIP account listing
- [x] Phone number (DID) listing
- [x] SMS sending
- [x] PBX info retrieval
- [x] JSON output support (`--json` flag)
- [x] Cobra-based CLI framework
- [x] SMS webhook daemon (`sms listen`)

### Phase 2
- [ ] DID number details (routing, forwarding)
- [ ] SMS history/logs
- [ ] Call rates and pricing
- [ ] Request callback endpoint
- [ ] Extension management
- [ ] Call webhook daemon (incoming call events)
- [ ] Log file output for daemon mode

### Phase 3
- [ ] Piping/STDIN support for bulk operations
- [ ] Call recording management
- [ ] PBX statistics and usage reports
- [ ] Config file support (.zadarma/config)
- [ ] Systemd service file for daemon mode

## Testing

Run unit tests:
```bash
go test ./internal/auth -v
```

Signature generation is tested against the known test credentials.

Project-wide tests:

```bash
go test ./...
```

Smoke script for quick validation (requires valid API credentials):

```bash
export ZADARMA_API_KEY=...
export ZADARMA_API_SECRET=...
./scripts/smoke.sh
```

The script builds the CLI and runs a subset of safe read-only commands. SMS send and extended statistics calls are
opt-in via environment flags to avoid accidental charges.

## Dependencies

Go 1.25.7 - uses only stdlib (no external dependencies).

## API Usage Notes & Decisions

This section captures practical details we validated while integrating with the Zadarma API.

### Endpoints implemented

- Balance: `GET /v1/info/balance/`
- SIP list: `GET /v1/sip/`
- SIP status: `GET /v1/sip/{id}/status/`
- SMS send: `POST /v1/sms/send/` (params in body; `caller_id` optional)
- SMS senders for numbers: `GET /v1/sms/senderid/?phones=...`
- PBX info: `GET /v1/pbx/`
- Webhook set: `POST /v1/pbx/webhooks/url/` (body: `url=<...>`)
- Webhook get: `GET /v1/pbx/webhooks/url/`
- Statistics: `GET /v1/statistics/` (supports `start`, `end`, `sip`, `cost_only`)
- Price lookup: `GET /v1/info/price/?number=...`
- Direct numbers:
  - Countries: `GET /v1/direct_numbers/countries/`
  - Country destinations: `GET /v1/direct_numbers/country/?country=...`
  - Number details: `GET /v1/direct_numbers/number/?type=...&number=...`
- DID list (user’s purchased numbers): `GET /v1/direct_numbers/`  ← Note: older docs/samples reference `/info/did/`,
  which returns 404 “Wrong method name” on production. We use `/direct_numbers/` per the official TypeScript client.

### Request building

- Non-GET parameters are sent in the body as `application/x-www-form-urlencoded`.
- GET parameters are signed and appended to the URL’s query.
- We removed implicit `format=json` from all requests. Commands decide presentation (text vs JSON) independent of
  transport.

### Logging

- Centralized logger under `internal/log` with `Debugf/Infof/Errorf` and a global `SetDebug()`.
- Enable with `-d/--debug` to print request method/URL, Authorization header, and the full raw response body for
  troubleshooting.

### Numbers formatting

- The official clients normalize phone numbers to digits only (strip non-digits). Our client currently accepts E.164 (
  e.g., `+1234567890`) as-is. If stricter normalization is required, we can add it in a backwards-compatible way.

### Sandbox vs production

- Our base URL targets production `https://api.zadarma.com/v1`.
- If sandbox support is required, we can add a constructor flag or environment switch to point to the sandbox host.

## Deployment to Forgejo

### Prerequisites
1. Create an empty repository on Forgejo: `https://forgejo.o-st.dev/zadarma/zadarma-cli`
2. SSH key configured for git access

### Push Commands
```bash
cd projects/zadarma-cli
git remote add origin ssh://git@forgejo.o-st.dev/zadarma/zadarma-cli.git
git branch -M main
git push -u origin main
```

## Next Steps

1. Create empty repository on Forgejo (https://forgejo.o-st.dev)
2. Push to Forgejo: `git push origin main`
3. Add more endpoint wrappers (rates, SMS, etc.)
4. Add integration tests with real API
5. Consider CLI framework upgrade (cobra) if feature set grows

## Commit History

### v0.1.0 - Initial Release
- HMAC-SHA1 authentication for Zadarma API
- Get account balance command
- Comprehensive auth signing implementation with tests
- Modular architecture: auth, client, CLI
- Zero external dependencies

## Release Conventions

### Versioning
- Use semantic versioning: `v0.0.1`, `v0.0.2`, `v0.1.0`, etc.
- Tags should be sequential and consistent
- **Never skip versions** - if v0.0.1 exists, next is v0.0.2, not v0.2.1

### Release Process
- **Releases are manual and intentional** - do NOT create a release after every change
- Release only when there's a meaningful milestone or user-facing change
- Before releasing:
  1. Update version in code if needed
  2. Ensure CI passes
  3. Update release notes in workflow
  4. Push tag manually

### CI/CD
- Tests run on every push to main
- Releases only trigger on explicit tag push (not on commits)
- Use `actions/forgejo-release@v2.11.1` for Forgejo releases
- **Always verify CI success after push or release** — check Forgejo Actions before moving on

### Files to Keep Out of Git
- `dist/` - build artifacts
- Binary files (zadarma, zadarma-test) - built by CI

## Development Workflow

### Code Organization & Separation of Concerns

- **Client Package** (`internal/client/`):
    - Core HTTP client and request handling in `client.go`
    - Simple utility methods (GetBalance, GetPrice) remain in `client.go`
    - Each major feature/command gets a dedicated file:
        - `sip.go` - SIP account methods
        - `sms.go` - SMS sending and senders
        - `direct_numbers.go` - Direct number/DID operations
        - `pbx.go` - PBX configuration methods
        - `statistics_client.go` - Statistics API wrapper

- **Commands Package** (`internal/commands/`):
    - CLI logic and command handlers
    - One file per command for clarity (e.g., `pbx.go`, `sms.go`, `statistics.go`)
    - Each command file handles:
        - Flag/argument parsing
        - Calling appropriate client methods
        - Output formatting (text vs JSON)

### Code Quality Standards

After **every change**, you must:

1. **Run linter**: Check code style and best practices
   ```bash
   golangci-lint run ./...
   ```
   Or if golangci-lint not available, use:
   ```bash
   go vet ./...
   ```

2. **Run formatter**: Ensure consistent code style
   ```bash
   go fmt ./...
   gofmt -s -w .
   ```

3. **Run tests**: Verify functionality and catch regressions
   ```bash
   go test ./...          # Run all tests
   go test ./... -v       # Verbose output
   go test -race ./...    # Check for race conditions
   ```

4. **Build**: Verify the project compiles
   ```bash
   go build -o zadarma-cli ./cmd/zadarma
   ```

### Pre-commit Checklist

Before committing code:

- [ ] All code is formatted (`go fmt ./...`)
- [ ] Linter passes (`go vet ./...` or `golangci-lint run ./...`)
- [ ] All tests pass (`go test ./...`)
- [ ] Project builds successfully (`go build ./...`)
- [ ] No new compile errors or warnings

### Development Example

```bash
# Make your code changes...
git add .

# Format code
go fmt ./...

# Check for issues
go vet ./...

# Run tests
go test ./...

# Build to verify
go build ./...

# If all pass, commit
git commit -m "description of changes"
```

## SMS Webhook Implementation

### Design Pattern: List vs Info Commands

The CLI follows a consistent pattern for all commands:

- `sms listen` - Listen for incoming SMS webhooks (requires webhook already configured)
- `sms set-webhook <URL>` - Register webhook URL and immediately start listening

### SMS Webhook Flow

The `sms` command group has the following structure:

**Prerequisites to receive SMS:**

1. Register the webhook URL with Zadarma
2. Enable SMS webhook notifications

**Two approaches:**

**Option A: Manual setup** (existing webhook)

```bash
zadarma-cli sms listen --port 8080
```

- Requires webhook already configured in Zadarma
- Exits with error if no webhook URL is configured
- Useful when webhook persists across sessions

**Option B: One-shot setup** (fresh webhook)

```bash
zadarma-cli sms set-webhook https://my-tunnel.loca.lt --port 8080
```

- Registers webhook URL
- Enables SMS notifications
- Immediately starts listening
- Combines webhook setup and listening

### Code Reuse Pattern

Both commands use the shared `startSMSListener()` helper function to:

- Set up HTTP listener on specified port
- Handle Zadarma verification requests (zd_echo)
- Parse and display incoming SMS events
- Support both text and JSON output formats

This avoids duplication and ensures consistent behavior across both listen paths.

### Implementation Details

**Webhook URL Validation:**

- Uses `url.ParseRequestURI()` to validate webhook URLs
- Ensures valid URL format before sending to API

**Response Parsing:**

- Normalizes API responses to handle both direct fields and nested data
- Checks status field for success/error conditions

**Listener HTTP Handler:**

- Responds to GET requests with zd_echo parameter (verification)
- Processes POST requests containing SMS events
- Returns HTTP 200 with "OK" on success

**Output Modes:**

- **Text mode (default):** Formatted table with FROM, TO, TEXT, TIME fields
- **JSON mode:** Raw JSON from webhook payload

