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
