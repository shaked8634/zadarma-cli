package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// GetPBXInfo fetches PBX configuration information. If pbxID or numbers are provided they are passed as query parameters.
// This method is resilient to two API shapes:
// 1) {"status":"success","data":{...}}
// 2) {"status":"success","pbx_id":...,"numbers":[...]}
func (c *Client) GetPBXInfo(pbxID, numbers string) (map[string]interface{}, error) {
	method := "/pbx/internal/"
	params := url.Values{}
	if pbxID != "" {
		params.Set("pbx_id", pbxID)
	}
	if numbers != "" {
		params.Set("numbers", numbers)
	}

	// Use a generic raw map to accept both response shapes
	var raw map[string]interface{}
	if err := c.Get(method, params, &raw); err != nil {
		return nil, err
	}

	// status check
	if st, _ := raw["status"].(string); st != "success" {
		msg := st
		if msg == "" {
			// try message field
			if m, _ := raw["message"].(string); m != "" {
				msg = m
			}
		}
		if msg == "" {
			return nil, fmt.Errorf("API error: unknown status")
		}
		return nil, fmt.Errorf("API error: %s", msg)
	}

	// Prefer nested data if present
	if d, ok := raw["data"].(map[string]interface{}); ok && d != nil {
		return d, nil
	}

	// Otherwise return the raw map without the status field
	// Make a shallow copy so we don't mutate the original
	out := map[string]interface{}{}
	for k, v := range raw {
		if k == "status" || k == "message" {
			continue
		}
		out[k] = v
	}

	// Ensure it's JSON-friendly (encode/decode) to normalize numbers to json types
	b, _ := json.Marshal(out)
	var norm map[string]interface{}
	_ = json.Unmarshal(b, &norm)

	return norm, nil
}

// SetWebhook sets a notification URL for events.
func (c *Client) SetWebhook(urlStr string) (map[string]interface{}, error) {
	// Validate webhook URL
	if err := ValidateWebhookURL(urlStr); err != nil {
		return nil, err
	}

	method := "/pbx/webhooks/url/"
	params := url.Values{}
	params.Set("url", urlStr)

	var resp struct {
		Status string                 `json:"status"`
		Data   map[string]interface{} `json:"data"`
	}

	// Send as form-encoded parameters (not JSON body)
	if err := c.Post(method, params, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// GetWebhooks returns the current notification URL.
func (c *Client) GetWebhooks() (map[string]interface{}, error) {
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

// SetWebhookHooks enables or disables webhook event types (e.g., sms).
// It posts to /pbx/webhooks/hooks with form-encoded parameters like sms=true
func (c *Client) SetWebhookHooks(enableSMS bool) (map[string]interface{}, error) {
	method := "/pbx/webhooks/hooks/"
	params := url.Values{}
	if enableSMS {
		params.Set("sms", "true")
	} else {
		params.Set("sms", "false")
	}

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

	// Normalize response
	b, _ := json.Marshal(resp.Data)
	var norm map[string]interface{}
	_ = json.Unmarshal(b, &norm)
	return norm, nil
}

// GetPBXInternalStatus fetches PBX internal status for a specific PBX ID.
// API: GET /v1/pbx/internal/{pbxId}/status/
func (c *Client) GetPBXInternalStatus(pbxID string) (map[string]interface{}, error) {
	method := "/pbx/internal/" + pbxID + "/status/"
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

// GetPBXInternalInfo fetches PBX internal info for a specific PBX ID.
// API: GET /v1/pbx/internal/{pbxId}/info/
func (c *Client) GetPBXInternalInfo(pbxID string) (map[string]interface{}, error) {
	method := "/pbx/internal/" + pbxID + "/info/"
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
