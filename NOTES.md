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

## Future Enhancements

### Phase 1 (Complete)
- [x] Get balance
- [x] SIP account listing
- [x] Phone number (DID) listing
- [x] SMS sending
- [x] PBX info retrieval
- [x] JSON output support (`--json` flag)
- [x] Cobra-based CLI framework

### Phase 2
- [ ] DID number details (routing, forwarding)
- [ ] SMS history/logs
- [ ] Call rates and pricing
- [ ] Request callback endpoint
- [ ] Extension management

### Phase 3
- [ ] Piping/STDIN support for bulk operations
- [ ] Call recording management
- [ ] PBX statistics and usage reports
- [ ] Config file support (.zadarma/config)

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
