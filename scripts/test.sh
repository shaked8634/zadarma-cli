#!/usr/bin/env bash
set -euo pipefail

export ZADARMA_API_KEY=acea883af4593167fe4a
export ZADARMA_API_SECRET=eec23d878ff592bb7a34

PBX_ID=649222


#cd .. && go build -o zadarma-cli ./cmd/zadarma
./zadarma-cli completion zsh > ~/.config/zsh/completions/_zadarma-cli
source ~/.config/zsh/completions/_zadarma-cli

echo "=== Testing webhook set (POST with body) ==="
./zadarma-cli -d webhook 'https://eight-dragons-thank.loca.lt'

echo ""
echo "=== Testing balance (GET) ==="
./zadarma-cli balance

echo ""
echo "=== Testing SIP list (GET) ==="
./zadarma-cli sip list

echo ""
echo "=== Testing SIP info (GET) ==="
./zadarma-cli sip info 649222

echo ""
echo "=== Testing phone list (GET) ==="
./zadarma-cli phone list

echo ""
echo "=== Testing phone number (GET with validation) ==="
./zadarma-cli phone number 972556620707

echo ""
echo "=== Testing phone countries list (GET) ==="
./zadarma-cli phone countries list

echo ""
echo "=== Testing phone country info (GET) ==="
./zadarma-cli phone country info IL

echo ""
echo "=== Testing pbx info (GET) ==="
./zadarma-cli pbx info

echo ""
echo "=== Testing webhook get (GET) ==="
./zadarma-cli webhook get

echo ""
echo "=== Testing webhook get (GET) ==="
./zadarma-cli sms senders --phones 972556620707

cd -