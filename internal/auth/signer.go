package auth

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/zadarma/zadarma-cli/internal/log"
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
	log.Debugf("paramsStr=%q", paramsStr)

	// Step 3: Calculate MD5 of params string
	paramsMD5 := s.md5Hex(paramsStr)
	log.Debugf("md5Hex=%q", paramsMD5)

	// Concatenate: method + paramsStr + md5(paramsStr)
	signString := method + paramsStr + paramsMD5
	logStr := signString
	if len(logStr) > 50 {
		logStr = logStr[:50]
	}
	log.Debugf("signString(first 50)=%q", logStr)

	// Step 4: HMAC-SHA1 with secret key
	hmacResult := s.hmacSHA1(signString)

	// Step 5: Base64 encode the hex-encoded HMAC (matches Python api.py)
	hashHex := hex.EncodeToString(hmacResult)
	logHex := hashHex
	if len(logHex) > 20 {
		logHex = logHex[:20]
	}
	log.Debugf("hashHex[:20]=%q", logHex)

	bts := []byte(hashHex)
	signature := base64.StdEncoding.EncodeToString(bts)
	log.Debugf("final_sig=%q", signature)

	return signature
}

// AuthHeader generates the full Authorization header value.
func (s *Signer) AuthHeader(method string, params url.Values) string {
	signature := s.Sign(method, params)
	return fmt.Sprintf("%s:%s", s.APIKey, signature)
}

// buildQueryString creates a query string from params using the same rules
// as the official clients (application/x-www-form-urlencoded):
// - keys are sorted
// - spaces are encoded as '+' (not '%20')
// Using url.Values.Encode() matches this behavior.
func (s *Signer) buildQueryString(params url.Values) string {
	if len(params) == 0 {
		return ""
	}
	// url.Values.Encode sorts keys and formats spaces as '+'
	return params.Encode()
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
