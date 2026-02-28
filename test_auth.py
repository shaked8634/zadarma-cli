#!/usr/bin/env python3
"""
Test Zadarma API authentication with both CORRECT and GO implementations.
"""

import hashlib
import hmac
import base64
import urllib.parse
import requests
import json

API_KEY = "acea883af4593167fe4a"
API_SECRET = "eec23d878ff592bb7a34"
METHOD = "/v1/info/balance/"
PARAMS = {"format": "json"}

# === CORRECT ZADARMA API WAY (base64(digest)) ===
def generate_signature_correct(api_key, api_secret, method, params):
    """Correct Zadarma API authentication."""
    # Sort and build query string
    sorted_params = sorted(params.items())
    request_line = urllib.parse.urlencode(sorted_params).replace('+', '%20')
    
    # MD5 of request_line
    md5_request_line = hashlib.md5(request_line.encode()).hexdigest()
    
    # Concatenate method + request_line + md5
    line = method + request_line + md5_request_line
    
    # HMAC-SHA1
    digest = hmac.new(api_secret.encode(), line.encode(), hashlib.sha1).digest()
    
    # Base64 encode the DIGEST (raw bytes)
    signature = base64.b64encode(digest).decode()
    
    return signature, request_line, md5_request_line, line, digest

# === GO IMPLEMENTATION WAY (base64(hex(digest))) ===
def generate_signature_go_way(api_key, api_secret, method, params):
    """Go implementation: base64(hex(digest))."""
    # Sort and build query string
    sorted_params = sorted(params.items())
    request_line = urllib.parse.urlencode(sorted_params).replace('+', '%20')
    
    # MD5 of request_line
    md5_request_line = hashlib.md5(request_line.encode()).hexdigest()
    
    # Concatenate method + request_line + md5
    line = method + request_line + md5_request_line
    
    # HMAC-SHA1
    digest = hmac.new(api_secret.encode(), line.encode(), hashlib.sha1).digest()
    
    # Hex encode the digest
    hex_digest = digest.hex()
    
    # Base64 encode the hex string
    signature = base64.b64encode(hex_digest.encode()).decode()
    
    return signature, request_line, md5_request_line, line, hex_digest

# Test both
print("=" * 80)
print("ZADARMA AUTH DEBUG - PYTHON TEST")
print("=" * 80)
print(f"API_KEY: {API_KEY}")
print(f"API_SECRET: {API_SECRET}")
print(f"METHOD: {METHOD}")
print(f"PARAMS: {PARAMS}")
print()

# CORRECT way
print("CORRECT WAY (base64 of digest):")
print("-" * 80)
sig_correct, req_line, md5_hash, line, digest = generate_signature_correct(API_KEY, API_SECRET, METHOD, PARAMS)
print(f"request_line (params): {req_line!r}")
print(f"md5_hash: {md5_hash!r}")
print(f"line (first 50): {line[:50]!r}")
print(f"digest (hex): {digest.hex()!r}")
print(f"digest[:20]: {digest.hex()[:20]!r}")
print(f"signature: {sig_correct!r}")
print(f"Authorization: {API_KEY}:{sig_correct}")
print()

# GO way (what current code does)
print("GO WAY (base64 of hex string):")
print("-" * 80)
sig_go, req_line_go, md5_hash_go, line_go, hex_digest = generate_signature_go_way(API_KEY, API_SECRET, METHOD, PARAMS)
print(f"request_line (params): {req_line_go!r}")
print(f"md5_hash: {md5_hash_go!r}")
print(f"line (first 50): {line_go[:50]!r}")
print(f"hex_digest[:20]: {hex_digest[:20]!r}")
print(f"signature: {sig_go!r}")
print(f"Authorization: {API_KEY}:{sig_go}")
print()

# Compare with actual GO output
GO_SIGNER_SIG = "ODEwNmNhYmFhMjBmODQ1ZTQ1OGRjOWZlMzQwNWNkZjQ0MmQ4NmQ0ZQ=="
print("GO ACTUAL OUTPUT:")
print("-" * 80)
print(f"signature: {GO_SIGNER_SIG!r}")
print(f"Authorization: {API_KEY}:{GO_SIGNER_SIG}")
print()

# Comparison
print("COMPARISON:")
print("-" * 80)
print(f"Correct == GO way?      {sig_correct == sig_go}")
print(f"Correct == GO actual?   {sig_correct == GO_SIGNER_SIG}")
print(f"GO way == GO actual?    {sig_go == GO_SIGNER_SIG}")
print()

# Try actual API calls
print("=" * 80)
print("API CALLS")
print("=" * 80)

# Test PROD
print("TESTING PRODUCTION URL:")
print("-" * 80)
base_urls = {
    "PROD": "https://api.zadarma.com",
    "SANDBOX": "https://sandbox.zadarma.com"
}

for env, base_url in base_urls.items():
    print(f"\n{env}: {base_url}")
    
    # Test with CORRECT signature
    url = f"{base_url}{METHOD}?format=json"
    headers = {
        "Authorization": f"{API_KEY}:{sig_correct}",
        "Content-Type": "application/x-www-form-urlencoded"
    }
    
    try:
        resp = requests.get(url, headers=headers, timeout=5)
        print(f"  [CORRECT] Status: {resp.status_code}")
        print(f"  Response: {resp.text[:100]}")
    except Exception as e:
        print(f"  [CORRECT] Error: {e}")
    
    # Test with GO way signature
    url = f"{base_url}{METHOD}?format=json"
    headers = {
        "Authorization": f"{API_KEY}:{sig_go}",
        "Content-Type": "application/x-www-form-urlencoded"
    }
    
    try:
        resp = requests.get(url, headers=headers, timeout=5)
        print(f"  [GO WAY] Status: {resp.status_code}")
        print(f"  Response: {resp.text[:100]}")
    except Exception as e:
        print(f"  [GO WAY] Error: {e}")
