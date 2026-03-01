# GetDirectNumbers Implementation Update Summary

## Overview

Successfully updated the `GetDirectNumbers()` method in the Zadarma CLI to support both listing all phone numbers and
fetching detailed information for specific numbers.

## Changes Made

### 1. Client Implementation (`internal/client/client.go`)

#### Updated `GetDirectNumbers()` Method

- **Signature**: Changed from `GetDirectNumbers()` to `GetDirectNumbers(numbers ...string)`
- **Behavior**:
    - **No arguments**: Fetches all DIDs via `GET /v1/direct_numbers/`
    - **With arguments**: Fetches detailed information for each specified number

#### New Helper Methods

- **`getAllDirectNumbers()`**: Internal method to fetch all phone numbers (implements the base endpoint)
- **`getDirectNumbersByList(numbers []string)`**: Internal method to fetch details for specific numbers
    - Optimized to call `getAllDirectNumbers()` only once
    - Builds a map of number→type for efficient lookup
    - Calls `/v1/direct_numbers/number/` endpoint for each requested number

### 2. Command Implementation (`internal/commands/phone_list.go`)

#### Updated Phone List Command

- **Command**: `zadarma-cli phone list [number...]`
- **Arguments**: Optional phone numbers in E.164 format (e.g., +14155555555)
- **Output Formats**:
    - **Text (default)**: Formatted table with columns: NUMBER, TYPE, STATUS, INFO
    - **JSON**: Full JSON response with `--output json` flag

#### New Functions

- **`handlePhoneList()`**: Main handler that calls `GetDirectNumbers()` and formats output
- **`printPhoneTable()`**: Renders tabular output with proper formatting
- **`getStringField()`**: Safely extracts string fields from response maps

### 3. Tests (`internal/client/client_test.go`)

#### Updated Tests

- **`TestGetDIDs()`**: Tests fetching all DIDs via `GetDirectNumbers()` with no arguments
- **`TestGetDIDsByNumber()`**: Tests fetching details for a single specific number
- **`TestGetDIDsByMultipleNumbers()`**: Tests fetching details for multiple numbers

All tests pass successfully ✓

### 4. Sample API Responses (`samples/`)

#### Files Created

1. **`v1_direct_numbers.json`**: Sample response from `GET /v1/direct_numbers/`
    - Contains array of 3 example numbers with basic info

2. **`v1_direct_numbers_single.json`**: Sample response from `GET /v1/direct_numbers/number/`
    - Contains detailed information for a single number
    - Includes extended fields like country_code, region, city, carrier, etc.

3. **`README_DIRECT_NUMBERS.md`**: Documentation for the sample files
    - Explains endpoint usage
    - Documents response structure
    - References official Zadarma API documentation

## Usage Examples

### List all phone numbers (text format)

```bash
zadarma-cli phone list
```

**Output:**

```
NUMBER          TYPE   STATUS   INFO
------          ----   ------   ----
+14155555555    voice  active   United States
+14155555556    fax    active   United States
+44207946824    voice  active   United Kingdom
```

### List all phone numbers (JSON format)

```bash
zadarma-cli phone list --output json
```

### Get details for a specific number

```bash
zadarma-cli phone list +14155555555
```

### Get details for multiple numbers

```bash
zadarma-cli phone list +14155555555 +14155555556
```

## API Endpoints Used

1. **`GET /v1/direct_numbers/`**
    - Returns all phone numbers (DIDs) owned by the account
    - Response: Array of phone number objects

2. **`GET /v1/direct_numbers/number/`**
    - Returns detailed information for a specific phone number
    - Parameters: `type` (required), `number` (required)
    - Response: Single phone number object with extended details

## Validation & Testing

✓ All existing tests pass
✓ All new tests pass
✓ Code builds successfully
✓ No compilation errors
✓ Implementation follows DRY principle with optimized API calls

## Architecture Notes

- **Single API call optimization**: When fetching multiple specific numbers, the implementation fetches the complete
  list once, then individually queries each number's details
- **Proper error handling**: Validates that requested numbers exist before attempting detailed queries
- **Consistent output formatting**: Uses tabwriter for proper text alignment and column formatting
- **JSON marshaling**: Preserves full API response data in JSON mode for advanced users
