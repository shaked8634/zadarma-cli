package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/zadarma/zadarma-cli/internal/auth"
	"github.com/zadarma/zadarma-cli/internal/log"
)

const (
	BaseURL    = "https://api.zadarma.com"
	SandboxURL = "https://api-sandbox.zadarma.com"
	APIVersion = "/v1"
)

// Client is the Zadarma API client.
type Client struct {
	baseURL string
	signer  *auth.Signer
	http    *http.Client
	debug   bool
}

// NewClient creates a new Zadarma API client.
// Use sandbox=true for the Zadarma sandbox API.
func NewClient(apiKey, apiSecret string, debug bool, sandbox bool) *Client {
	// configure global logger debug flag based on client setting
	log.SetDebug(debug)
	baseURL := BaseURL
	if sandbox {
		baseURL = SandboxURL
	}
	return &Client{
		baseURL: baseURL + APIVersion,
		signer:  auth.NewSigner(apiKey, apiSecret),
		http:    &http.Client{},
		debug:   debug,
	}
}

// APIResponse is the base structure for all Zadarma API responses.
type APIResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// GetBalance fetches the user's account balance.
func (c *Client) GetBalance() (balance interface{}, currency string, err error) {
	method := "/info/balance/"
	params := url.Values{}

	var resp struct {
		Status   string  `json:"status"`
		Balance  float64 `json:"balance"`
		Currency string  `json:"currency"`
	}

	if err := c.Get(method, params, &resp); err != nil {
		return nil, "", err
	}

	if resp.Status != "success" {
		return nil, "", fmt.Errorf("API error: %s", resp.Status)
	}

	return resp.Balance, resp.Currency, nil
}

// GetPrice returns the price information for a call to the given number.
// API: GET /v1/info/price/?number=<phone>
func (c *Client) GetPrice(number string) (map[string]interface{}, error) {
	method := "/info/price/"
	params := url.Values{}
	params.Set("number", number)

	var resp struct {
		Status string                 `json:"status"`
		Data   map[string]interface{} `json:"data"`
	}

	if err := c.Get(method, params, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	return resp.Data, nil
}

// ...methods moved to dedicated files: sip.go, sms.go, direct_numbers.go, pbx.go, statistics_client.go...

// Get performs a GET request to the API.
func (c *Client) Get(method string, params url.Values, result interface{}) error {
	return c.request("GET", method, params, nil, result)
}

// Post performs a POST request to the API.
func (c *Client) Post(method string, params url.Values, body io.Reader, result interface{}) error {
	return c.request("POST", method, params, body, result)
}

// PostJSON sends a JSON-encoded body but uses an explicit empty signingParams (per API rules)
// so the signature is calculated over empty params while the body is JSON.
func (c *Client) PostJSON(method string, jsonBody interface{}, result interface{}) error {
	b, err := json.Marshal(jsonBody)
	if err != nil {
		return fmt.Errorf("failed to marshal json body: %w", err)
	}
	// For PostJSON, we DO NOT include these fields in signingParams; signature is based on empty params
	return c.doRequest(http.MethodPost, method, nil, bytes.NewReader(b), result, url.Values{}, "application/json")
}

// request performs an HTTP request with proper Zadarma authentication using params for signing.
func (c *Client) request(httpMethod, apiMethod string, params url.Values, body io.Reader, result interface{}) error {
	// By default signingParams == params, and content type for POST/PUT is form-encoded when appropriate
	contentType := ""
	if httpMethod == http.MethodPost || httpMethod == http.MethodPut {
		// If body is nil and params exist, we'll send form-encoded body and set content-type accordingly in doRequest
		contentType = "application/x-www-form-urlencoded"
	}
	return c.doRequest(httpMethod, apiMethod, params, body, result, params, contentType)
}

// doRequest is the underlying HTTP request implementation. signingParams are used to compute the Auth header
// (they do NOT have to match the actual HTTP body); params is used to build the URL query string for GET requests
// and optionally to form-encode the body for POST when body==nil.
func (c *Client) doRequest(httpMethod, apiMethod string, params url.Values, body io.Reader, result interface{}, signingParams url.Values, contentType string) error {
	// Initialize params if nil
	if params == nil {
		params = url.Values{}
	}
	if signingParams == nil {
		signingParams = url.Values{}
	}

	// Build URL and (possibly) request body depending on HTTP method
	fullURL := c.baseURL + apiMethod
	reqBody := body

	if httpMethod == http.MethodGet {
		if len(params) > 0 {
			fullURL = fullURL + "?" + params.Encode()
		}
	} else {
		// For non-GET methods send params in the body as x-www-form-urlencoded unless a body is supplied
		if reqBody == nil && len(params) > 0 {
			encoded := params.Encode()
			reqBody = strings.NewReader(encoded)
			// If contentType not explicitly set, default to form-encoded
			if contentType == "" {
				contentType = "application/x-www-form-urlencoded"
			}
		}
	}

	log.Debugf("Request: %s %s", httpMethod, fullURL)
	log.Debugf("paramsStr=%q", params.Encode())
	log.Debugf("signingParamsStr=%q", signingParams.Encode())

	// Create request
	req, err := http.NewRequest(httpMethod, fullURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Log the actual body content for debugging (whether derived from params or provided body)
	if reqBody != nil {
		// If we created the body from params (body == nil), log params.Encode(); otherwise we log that body was provided
		if body == nil {
			log.Debugf("Request body (from params): %s", params.Encode())
		} else {
			log.Debugf("Request body provided (reader)")
		}
	}

	// Sign request with full path including /v1
	signingPath := APIVersion + apiMethod
	authHeader := c.signer.AuthHeader(signingPath, signingParams)
	// Correct header name is 'Authorization'
	req.Header.Set("Authorization", authHeader)

	// Set Content-Type if present
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	log.Debugf("Authorization: %s", authHeader)
	if reqBody != nil {
		// Log request body if present
		if body != nil {
			// If body was provided, we already logged params above
			log.Debugf("Request body: %v", body)
		}
	}

	// Execute request
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Debugf("Error closing response body: %v", err)
		}
	}()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	log.Debugf("Response: HTTP %d (%d bytes). Body: %s", resp.StatusCode, len(respBody), string(respBody))

	// Check HTTP status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Debugf("Error response: %s", string(respBody))
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse JSON response
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	log.Debugf("Parsed response successfully")
	return nil
}
