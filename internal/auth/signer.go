package auth

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// Signer handles Zadarma API authentication via HMAC-SHA1 signatures.
type Signer struct {
	APIKey    string
	APISecret string
}

// NewSigner creates a new API signer with the given credentials.
func NewSigner(apiKey, apiSecret string) *Signer {
	return &Signer{
		APIKey:    apiKey,
		APISecret: apiSecret,
	}
}

// Sign generates an HMAC-SHA1 signature for the given method and parameters.
// Algorithm (from Zadarma docs):
// 1. Sort params alphabetically
// 2. Build query string: param1=value1&param2=value2
// 3. Concatenate: method + paramsStr + md5(paramsStr)
// 4. HMAC-SHA1 with secret key
// 5. Base64 encode
func (s *Signer) Sign(method string, params url.Values) string {
	// Step 1 & 2: Build alphabetically-sorted query string
	paramsStr := s.buildQueryString(params)

	// Step 3: Calculate MD5 of params string
	paramsMD5 := s.md5Hex(paramsStr)

	// Concatenate: method + paramsStr + md5(paramsStr)
	signString := method + paramsStr + paramsMD5

	// Step 4: HMAC-SHA1 with secret key
	hmacResult := s.hmacSHA1(signString)

	// Step 5: Base64 encode
	signature := base64.StdEncoding.EncodeToString(hmacResult)

	return signature
}

// AuthHeader generates the full Authorization header value.
func (s *Signer) AuthHeader(method string, params url.Values) string {
	signature := s.Sign(method, params)
	return fmt.Sprintf("%s:%s", s.APIKey, signature)
}

// buildQueryString creates an alphabetically-sorted query string from params.
func (s *Signer) buildQueryString(params url.Values) string {
	if len(params) == 0 {
		return ""
	}

	// Get sorted keys
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build query string with sorted params
	var parts []string
	for _, k := range keys {
		// Use the first value for each key
		parts = append(parts, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(params.Get(k))))
	}

	return strings.Join(parts, "&")
}

// md5Hex calculates MD5 hash and returns hex string.
func (s *Signer) md5Hex(data string) string {
	h := md5.Sum([]byte(data))
	return hex.EncodeToString(h[:])
}

// hmacSHA1 computes HMAC-SHA1 with the secret key.
func (s *Signer) hmacSHA1(data string) []byte {
	h := hmac.New(sha1.New, []byte(s.APISecret))
	h.Write([]byte(data))
	return h.Sum(nil)
}
