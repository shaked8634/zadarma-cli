#!/bin/bash
# Comprehensive test of the new phone command consolidation

set -e

export ZADARMA_API_KEY=acea883af4593167fe4a
export ZADARMA_API_SECRET=eec23d878ff592bb7a34

cd /home/tam/Projects/zadarma-cli

echo "=== Testing Zadarma CLI Phone Command ==="
echo ""

echo "1. Testing: zadarma-cli phone list (text format)"
./zadarma-cli phone list
echo ""

echo "2. Testing: zadarma-cli phone list (JSON format)"
./zadarma-cli phone list --output json | head -20
echo "..."
echo ""

echo "3. Testing: zadarma-cli phone list <number> (specific number)"
./zadarma-cli phone list 123456789012
echo ""

echo "4. Testing: zadarma-cli phone list <multiple numbers>"
./zadarma-cli phone list 123456789012 19293091254 --output json | head -15
echo "..."
echo ""

echo "5. Testing: zadarma-cli balance (verify other commands still work)"
./zadarma-cli balance
echo ""

echo "6. Testing: zadarma-cli phone countries"
./zadarma-cli phone countries 2>/dev/null | head -10
echo ""

echo "=== All tests completed successfully ==="

