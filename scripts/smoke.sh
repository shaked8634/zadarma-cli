#!/usr/bin/env bash
set -euo pipefail

# Simple smoke test runner for zadarma-cli.
# - Runs unit tests
# - Builds the CLI
# - Optionally runs a subset of live API commands if credentials are present
#
# Env vars:
#   ZADARMA_API_KEY / ZADARMA_API_SECRET   Required for live API calls
#   JSON_FLAG        Default: --json
#   DEBUG_FLAG       Default: -d
#   SMOKE_DIRECT_COUNTRY           e.g., US
#   SMOKE_DIRECT_NUMBER_TYPE       e.g., mobile|landline
#   SMOKE_DIRECT_NUMBER            e.g., 14155551234
#   SMOKE_SMS_NUMBER               Opt-in; number to send SMS to
#   SMOKE_SMS_TEXT                 Opt-in; text to send
#   SMOKE_SMS_SENDER               Optional sender (must be confirmed)
#   SMOKE_STATISTICS_START         Optional start date "YYYY-MM-DD HH:MM:SS"
#   SMOKE_STATISTICS_END           Optional end date "YYYY-MM-DD HH:MM:SS"

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$ROOT"

echo "[smoke] Running unit tests..."
go test ./... -count=1

echo "[smoke] Building CLI..."
go build -o zadarma-cli ./cmd/zadarma

if [[ -z "${ZADARMA_API_KEY:-}" || -z "${ZADARMA_API_SECRET:-}" ]]; then
  echo "[smoke] Skipping live API run: ZADARMA_API_KEY and/or ZADARMA_API_SECRET not set."
  exit 0
fi

JSON_FLAG=${JSON_FLAG:---json}
DEBUG_FLAG=${DEBUG_FLAG:--d}

run() { echo "+ $*"; "$@"; }

echo "[smoke] Running read-only commands against live API..."
run ./zadarma-cli "$DEBUG_FLAG" "$JSON_FLAG" balance
run ./zadarma-cli "$DEBUG_FLAG" "$JSON_FLAG" sip list
run ./zadarma-cli "$DEBUG_FLAG" "$JSON_FLAG" did list
run ./zadarma-cli "$DEBUG_FLAG" "$JSON_FLAG" pbx info
run ./zadarma-cli "$DEBUG_FLAG" "$JSON_FLAG" direct countries

if [[ -n "${SMOKE_DIRECT_COUNTRY:-}" ]]; then
  run ./zadarma-cli "$DEBUG_FLAG" "$JSON_FLAG" direct country "$SMOKE_DIRECT_COUNTRY"
fi

if [[ -n "${SMOKE_DIRECT_NUMBER_TYPE:-}" && -n "${SMOKE_DIRECT_NUMBER:-}" ]]; then
  run ./zadarma-cli "$DEBUG_FLAG" "$JSON_FLAG" direct number "$SMOKE_DIRECT_NUMBER_TYPE" "$SMOKE_DIRECT_NUMBER"
fi

# SMS sending is opt-in to avoid accidental charges
if [[ -n "${SMOKE_SMS_NUMBER:-}" && -n "${SMOKE_SMS_TEXT:-}" ]]; then
  if [[ -n "${SMOKE_SMS_SENDER:-}" ]]; then
    run ./zadarma-cli "$DEBUG_FLAG" "$JSON_FLAG" sms send --phone "$SMOKE_SMS_NUMBER" --message "$SMOKE_SMS_TEXT" --sender "$SMOKE_SMS_SENDER"
  else
    run ./zadarma-cli "$DEBUG_FLAG" "$JSON_FLAG" sms send --phone "$SMOKE_SMS_NUMBER" --message "$SMOKE_SMS_TEXT"
  fi
fi

# Statistics optional
if [[ -n "${SMOKE_STATISTICS_START:-}" || -n "${SMOKE_STATISTICS_END:-}" ]]; then
  ARGS=(statistics)
  [[ -n "${SMOKE_STATISTICS_START:-}" ]] && ARGS+=("--start" "$SMOKE_STATISTICS_START")
  [[ -n "${SMOKE_STATISTICS_END:-}" ]] && ARGS+=("--end" "$SMOKE_STATISTICS_END")
  run ./zadarma-cli "$DEBUG_FLAG" "$JSON_FLAG" "${ARGS[@]}"
fi

echo "[smoke] Done."
