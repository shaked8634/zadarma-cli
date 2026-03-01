# GetDirectNumbers Implementation - Quick Reference

## What Changed

The `GetDirectNumbers()` method has been updated to support two distinct API workflows:

### 1. **List All Numbers** (No Arguments)

```go
dids, err := c.GetDirectNumbers() // No arguments
```

- API Endpoint: `GET /v1/direct_numbers/`
- Returns: Array of all phone numbers with basic info

### 2. **Get Specific Number Details** (With Arguments)

```go
dids, err := c.GetDirectNumbers("+14155555555") // Single number
dids, err := c.GetDirectNumbers("+14155555555", "+14155556") // Multiple numbers
```

- API Endpoint: `GET /v1/direct_numbers/number/?type=voice&number=+14155555555`
- Returns: Array of detailed info for each requested number

## CLI Usage

### All Numbers (Text)

```bash
zadarma-cli phone list
```

### All Numbers (JSON)

```bash
zadarma-cli phone list --output json
```

### Specific Number(s)

```bash
zadarma-cli phone list 972556620707
zadarma-cli phone list 972556620707 19293091254
```

## Key Implementation Details

✓ **Optimized**: Fetches the complete list once when requesting specific numbers  
✓ **Validating**: Ensures requested numbers exist before querying details  
✓ **Error Handling**: Clear error messages if numbers aren't found  
✓ **Logging**: Includes debug logging via `log.Debugf()`  
✓ **Tested**: Full test coverage for all scenarios

## Files Modified

1. `internal/client/client.go` - Updated GetDirectNumbers and added helpers
2. `internal/commands/did.go` - Updated command handler and output formatting
3. `internal/client/client_test.go` - Added new test cases

## Files Created

1. `samples/v1_direct_numbers.json` - Sample API response for list all
2. `samples/v1_direct_numbers_single.json` - Sample API response for single number
3. `samples/README_DIRECT_NUMBERS.md` - API documentation and usage guide

## Test Results

All tests passing:

- ✓ TestGetDIDs (fetch all)
- ✓ TestGetDIDsByNumber (fetch single)
- ✓ TestGetDIDsByMultipleNumbers (fetch multiple)
- ✓ All other existing tests still pass

## Backward Compatibility

The change is **fully backward compatible**:

- `GetDirectNumbers()` with no arguments works exactly as before
- Existing code continues to work without modification
- New functionality is additive through variadic parameters
