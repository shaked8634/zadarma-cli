# User Requirements - Zadarma CLI Development ✅ COMPLETED

## API & Testing Rules

1. ✅ **Only save real API responses** - Sample files in `samples/` contain actual API responses from live calls
    - `samples/v1_direct_numbers.json` - Real response from GET /v1/direct_numbers/
    - `samples/v1_direct_numbers_single.json` - Real response for specific number query

2. ✅ **Use smoke.sh for testing** - Available at `scripts/smoke.sh` for live API testing
    - Supports real credentials via environment variables
    - Tests read-only and optional write operations

3. ✅ **Remove --json switch** - Removed the `--json` persistent flag
    - Use `--output json` instead for JSON output format
    - Updated wantsJSON() function to only check --output flag

4. ✅ **Compile to zadarma-cli binary** - Build output available as `zadarma-cli` in project root
    - Run: `go build -o zadarma-cli ./cmd/zadarma`
    - Binary is ready for immediate use

## Command Structure Changes

5. ✅ **Remove 'did' command** - Consolidated DID functionality into 'phone' command
    - Old `internal/commands/did.go` removed
    - All DID functionality now in `internal/commands/phone.go`

6. ✅ **Rename 'direct' to 'phone'** - Main phone/number management command
    - Old `internal/commands/direct.go` removed
    - Virtual number exploration also in `internal/commands/phone.go`

## Final Command Structure

```bash
zadarma-cli phone list                      # List all owned phone numbers
zadarma-cli phone list <number>             # Get details for specific number
zadarma-cli phone list <num1> <num2> ...    # Get details for multiple numbers
zadarma-cli phone countries                 # List available countries
zadarma-cli phone country <code>            # List destinations for country
zadarma-cli phone number <type> <number>    # Get virtual number info
```

## Implementation Details

- Updated `cmd/zadarma/main.go` - Removed DID and Direct registrations, added Phone
- Created `internal/commands/phone.go` - Unified phone management command
- Updated `internal/client/client.go` - Fixed API response handling (uses 'info' not 'data')
- Updated `internal/commands/common.go` - Removed --json flag logic
- All tests updated and passing
- Real API responses saved as samples
- Binary compiled and tested with live API

## Verification

Test with:

```bash
export ZADARMA_API_KEY=<your_key>
export ZADARMA_API_SECRET=<your_secret>
./zadarma-cli phone list --output json
```

All requirements have been successfully implemented! ✅
