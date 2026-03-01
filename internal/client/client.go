package client

import (
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
	// configure global logger debug flag based on client setting
	log.SetDebug(debug)
	return &Client{
		baseURL: BaseURL + APIVersion,
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

// GetSIPStatus fetches the status of a specific SIP account.
func (c *Client) GetSIPStatus(id string) (isOnline bool, err error) {
	method := "/sip/" + id + "/status/"
	params := url.Values{}

	var resp struct {
		Status   string `json:"status"`
		SIP      string `json:"sip"`
		IsOnline string `json:"is_online"`
	}

	if err := c.Get(method, params, &resp); err != nil {
		return false, err
	}

	if resp.Status != "success" {
		return false, fmt.Errorf("API error: %s", resp.Status)
	}

	return resp.IsOnline == "true", nil
}

// GetDirectNumbers fetches all phone numbers (DIDs) owned by the user.
// API (per official TS client): GET /v1/direct_numbers/
func (c *Client) GetDirectNumbers() ([]map[string]interface{}, error) {
	method := "/direct_numbers/"
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

// GetDirectCountries lists available direct number countries.
func (c *Client) GetDirectCountries() ([]map[string]interface{}, error) {
	method := "/direct_numbers/countries/"
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

// GetDirectCountry lists destinations for a specific country.
func (c *Client) GetDirectCountry(country string) ([]map[string]interface{}, error) {
	method := "/direct_numbers/country/"
	params := url.Values{}
	params.Set("country", country)

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

// GetDirectNumber returns information about a direct number.
func (c *Client) GetDirectNumber(type_, number string) (map[string]interface{}, error) {
	method := "/direct_numbers/number/"
	params := url.Values{}
	params.Set("type", type_)
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

// SendSMS sends an SMS message.
func (c *Client) SendSMS(phoneNumber, message, sender string) (map[string]interface{}, error) {
	method := "/sms/send/"
	params := url.Values{}
	params.Set("number", phoneNumber)
	params.Set("message", message)
	if sender != "" {
		params.Set("caller_id", sender)
	}

	// Use a generic map to accommodate different possible API shapes
	var raw map[string]any
	if err := c.Post(method, params, nil, &raw); err != nil {
		return nil, err
	}

	status, _ := raw["status"].(string)
	if status != "success" {
		if status == "" {
			status = "unknown"
		}
		return nil, fmt.Errorf("API error: %s", status)
	}

	// Prefer nested data map if present; otherwise, normalize common top-level fields
	out := map[string]any{}
	if d, ok := raw["data"].(map[string]any); ok && d != nil {
		out = d
	}
	// Normalize id field
	if _, ok := out["id"]; !ok || out["id"] == nil {
		if v, ok := raw["id"]; ok {
			out["id"] = v
		} else if v, ok := raw["message_id"]; ok {
			out["id"] = v
		} else if v, ok := raw["sms_id"]; ok {
			out["id"] = v
		}
	}
	// Normalize status field for SMS send result if present at alternative keys
	if _, ok := out["status"]; !ok || out["status"] == nil {
		if v, ok := raw["sms_status"]; ok {
			out["status"] = v
		} else if v, ok := raw["message_status"]; ok {
			out["status"] = v
		} else if v, ok := raw["status"]; ok { // fallback to overall status
			out["status"] = v
		}
	}

	return out, nil
}

// GetSMSSenders returns the list of valid SMS senders to a given number.
func (c *Client) GetSMSSenders(phones string) ([]map[string]interface{}, error) {
	method := "/sms/senderid/"
	params := url.Values{}
	params.Set("phones", phones)

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
	// Initialize params if nil
	if params == nil {
		params = url.Values{}
	}

	// Do not force format=json globally; callers/commands may add it explicitly if needed.

	// Build URL and (possibly) request body depending on HTTP method
	fullURL := c.baseURL + apiMethod
	reqBody := body
	if httpMethod == http.MethodGet {
		if len(params) > 0 {
			fullURL = fullURL + "?" + params.Encode()
		}
	} else {
		// For non-GET methods send params in the body as x-www-form-urlencoded
		if reqBody == nil && len(params) > 0 {
			reqBody = strings.NewReader(params.Encode())
		}
	}

	log.Debugf("Request: %s %s", httpMethod, fullURL)

	// Create request
	req, err := http.NewRequest(httpMethod, fullURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Sign request with full path including /v1 (format=json is already in params)
	// Extract version from baseURL to reconstruct full path for signing
	signingPath := APIVersion + apiMethod
	authHeader := c.signer.AuthHeader(signingPath, params)
	req.Header.Set("Authorization", authHeader)
	// Per docs, POST and PUT must specify Content-Type
	if httpMethod == http.MethodPost || httpMethod == http.MethodPut {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	log.Debugf("Authorization: %s", authHeader)

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

	log.Debugf("Response: HTTP %d (%d bytes)", resp.StatusCode, len(respBody))
	// Per user request, print the full response body
	log.Debugf("Response body: %s", string(respBody))

	// Check HTTP status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Debugf("Error response: %s", string(respBody))
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse JSON response
	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	log.Debugf("Parsed response successfully")
	return nil
}
