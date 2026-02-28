package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/zadarma/zadarma-cli/internal/auth"
)

const (
	BaseURL    = "https://api.zadarma.com"
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
func NewClient(apiKey, apiSecret string, debug bool) *Client {
	return &Client{
		baseURL: BaseURL + APIVersion,
		signer:  auth.NewSigner(apiKey, apiSecret),
		http:    &http.Client{},
		debug:   debug,
	}
}

// debugPrint prints debug messages if debug mode is enabled.
func (c *Client) debugPrint(format string, args ...interface{}) {
	if c.debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
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

// GetSIPs fetches all SIP accounts.
func (c *Client) GetSIPs() ([]map[string]interface{}, error) {
	method := "/sip/"
	params := url.Values{}

	var resp struct {
		Status string                   `json:"status"`
		Data   []map[string]interface{} `json:"data"`
	}

	if err := c.Get(method, params, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	return resp.Data, nil
}

// GetDIDs fetches all phone numbers (DIDs).
func (c *Client) GetDIDs() ([]map[string]interface{}, error) {
	method := "/info/did/"
	params := url.Values{}

	var resp struct {
		Status string                   `json:"status"`
		Data   []map[string]interface{} `json:"data"`
	}

	if err := c.Get(method, params, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	return resp.Data, nil
}

// SendSMS sends an SMS message.
func (c *Client) SendSMS(phoneNumber, message string) (map[string]interface{}, error) {
	method := "/sms/"
	params := url.Values{}
	params.Set("number", phoneNumber)
	params.Set("message", message)

	var resp struct {
		Status string                 `json:"status"`
		Data   map[string]interface{} `json:"data"`
	}

	if err := c.Post(method, params, nil, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	return resp.Data, nil
}

// GetPBXInfo fetches PBX configuration information.
func (c *Client) GetPBXInfo() (map[string]interface{}, error) {
	method := "/pbx/"
	params := url.Values{}

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

// SetWebhook sets a notification URL for events.
func (c *Client) SetWebhook(urlStr string) (map[string]interface{}, error) {
	method := "/pbx/webhooks/url/"
	params := url.Values{}
	params.Set("url", urlStr)

	var resp struct {
		Status string                 `json:"status"`
		Data   map[string]interface{} `json:"data"`
	}

	if err := c.Post(method, params, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// GetWebhook returns the current notification URL.
func (c *Client) GetWebhook() (map[string]interface{}, error) {
	method := "/pbx/webhooks/url/"
	params := url.Values{}

	var resp struct {
		Status string                 `json:"status"`
		Data   map[string]interface{} `json:"data"`
	}

	if err := c.Get(method, params, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// Get performs a GET request to the API.
func (c *Client) Get(method string, params url.Values, result interface{}) error {
	return c.request("GET", method, params, nil, result)
}

// Post performs a POST request to the API.
func (c *Client) Post(method string, params url.Values, body io.Reader, result interface{}) error {
	return c.request("POST", method, params, body, result)
}

// request performs an HTTP request with proper Zadarma authentication.
func (c *Client) request(httpMethod, apiMethod string, params url.Values, body io.Reader, result interface{}) error {
	// Build full URL
	fullURL := c.baseURL + apiMethod
	if len(params) > 0 {
		fullURL = fullURL + "?" + params.Encode()
	}

	c.debugPrint("Request: %s %s", httpMethod, fullURL)

	// Create request
	req, err := http.NewRequest(httpMethod, fullURL, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Sign request
	authHeader := c.signer.AuthHeader(apiMethod, params)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	c.debugPrint("Authorization: %s", authHeader[:30]+"...")

	// Execute request
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	c.debugPrint("Response: HTTP %d (%d bytes)", resp.StatusCode, len(respBody))

	// Check HTTP status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.debugPrint("Error response: %s", string(respBody))
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse JSON response
	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	c.debugPrint("Parsed response successfully")
	return nil
}
