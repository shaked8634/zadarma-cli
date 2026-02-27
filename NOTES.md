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

Currently supports **text output**:
```
Balance: 123.45 USD
```

JSON output can be added as a `--json` flag in the future.

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

### Testing the Balance Endpoint

Use this curl command to test (signature shown for your API key):
```bash
curl -X GET 'https://api.zadarma.com/v1/info/balance/' \
  -H 'Authorization: 1c64dee7ee7638e3a507:rhSjhy87PyCS8HamPv6b1t199EA='
```

To generate signatures for other endpoints:
1. Prepare your params as `key=value&key2=value2` (alphabetically sorted)
2. Calculate `md5(params_string)`
3. Build: `method + params + md5`
4. HMAC-SHA1 with secret
5. Base64 encode

The `Signer` package automates this.

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

## Dependencies

Go 1.25.7 - uses only stdlib (no external dependencies).

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

### Files to Keep Out of Git
- `dist/` - build artifacts
- Binary files (zadarma, zadarma-test) - built by CI
