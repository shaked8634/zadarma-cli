#!/usr/bin/env python3
"""
Test Zadarma API authentication - checking if method path is the issue.
"""

import hashlib
import hmac
import base64
import urllib.parse
import requests

API_KEY = "acea883af4593167fe4a"
API_SECRET = "eec23d878ff592bb7a34"
PARAMS = {"format": "json"}

def generate_signature_go_way(api_key, api_secret, method, params):
    """Go implementation: base64(hex(HMAC-SHA1))."""
    sorted_params = sorted(params.items())
    request_line = urllib.parse.urlencode(sorted_params).replace('+', '%20')
    
    md5_request_line = hashlib.md5(request_line.encode()).hexdigest()
    line = method + request_line + md5_request_line
    
    digest = hmac.new(api_secret.encode(), line.encode(), hashlib.sha1).digest()
    hex_digest = digest.hex()
    signature = base64.b64encode(hex_digest.encode()).decode()
    
    return signature, request_line, md5_request_line, line, hex_digest, digest

print("=" * 80)
print("METHOD PATH COMPARISON")
print("=" * 80)

methods = {
    "WITH /v1": "/v1/info/balance/",
    "WITHOUT /v1": "/info/balance/"
}

go_actual_sig = "ODEwNmNhYmFhMjBmODQ1ZTQ1OGRjOWZlMzQwNWNkZjQ0MmQ4NmQ0ZQ=="
go_actual_hex = base64.b64decode(go_actual_sig).decode()

print(f"\nGO actual signature: {go_actual_sig}")
print(f"GO actual hex (decoded): {go_actual_hex}")
print()

for label, method in methods.items():
    print(f"\n{label}: {method}")
    print("-" * 80)
    sig, req_line, md5_hash, line, hex_digest, digest = generate_signature_go_way(API_KEY, API_SECRET, method, PARAMS)
    
    print(f"request_line: {req_line!r}")
    print(f"md5_hash: {md5_hash!r}")
    print(f"line (first 50): {line[:50]!r}")
    print(f"hex_digest: {hex_digest!r}")
    print(f"signature: {sig!r}")
    print(f"Match GO actual? {hex_digest == go_actual_hex}")
    
    # Try API call
    base_url = "https://api.zadarma.com"
    url = f"{base_url}{method}?format=json"
    headers = {
        "Authorization": f"{API_KEY}:{sig}",
        "Content-Type": "application/x-www-form-urlencoded"
    }
    
    try:
        resp = requests.get(url, headers=headers, timeout=5)
        print(f"API Response: HTTP {resp.status_code}")
        if resp.status_code == 200:
            data = resp.json()
            print(f"  Balance: {data.get('balance')} {data.get('currency')}")
        else:
            print(f"  Error: {resp.text[:50]}")
    except Exception as e:
        print(f"API Error: {e}")
