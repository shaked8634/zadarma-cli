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

	// Step 3: Calculate MD5 of params string
	paramsMD5 := s.md5Hex(paramsStr)

	// Concatenate: method + paramsStr + md5(paramsStr)
	signString := method + paramsStr + paramsMD5

	log.Debugf("fullSignString=%q", signString)

	// Step 4: HMAC-SHA1 with secret key
	hmacResult := s.hmacSHA1(signString)

	// Step 5 variants for debugging
	// hex->base64 (original approach used by this client)
	hashHex := hex.EncodeToString(hmacResult)
	hexBts := []byte(hashHex)
	hexBase64 := base64.StdEncoding.EncodeToString(hexBts)
	log.Debugf("signature(hex->base64)=%q", hexBase64)

	// raw-hmac -> base64 (alternative)
	rawBase64 := base64.StdEncoding.EncodeToString(hmacResult)
	log.Debugf("signature(raw-hmac->base64)=%q", rawBase64)

	// Return hex->base64 to preserve compatibility with GET requests and existing behavior.
	return hexBase64
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
