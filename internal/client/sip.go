package client

import (
	"fmt"
	"net/url"
)

// GetSIPs fetches all SIP accounts.
func (c *Client) GetSIPs() ([]map[string]interface{}, error) {
	method := "/sip/"
	params := url.Values{}

	var resp struct {
		Status string                   `json:"status"`
		SIPs   []map[string]interface{} `json:"sips"`
	}

	if err := c.Get(method, params, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	return resp.SIPs, nil
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

// SetSipCallerID sets the caller ID for a specific SIP account.
func (c *Client) SetSipCallerID(id, number string) (map[string]interface{}, error) {
	method := "/sip/" + id + "/callerid/"
	params := url.Values{}
	params.Set("caller_id", number)

	var raw map[string]interface{}
	if err := c.Post(method, params, nil, &raw); err != nil {
		return nil, err
	}

	if st, _ := raw["status"].(string); st != "success" {
		msg := st
		if m, _ := raw["message"].(string); m != "" {
			msg = m
		}
		if msg == "" {
			return nil, fmt.Errorf("API error: unknown status")
		}
		return nil, fmt.Errorf("API error: %s", msg)
	}

	// Return everything except status
	out := map[string]interface{}{}
	for k, v := range raw {
		if k == "status" || k == "message" {
			continue
		}
		out[k] = v
	}

	return out, nil
}
