package client

import (
	"fmt"
	"net/url"

	"github.com/zadarma/zadarma-cli/internal/log"
)

// GetDirectNumbers fetches phone numbers (DIDs) owned by the user.
// If numbers is empty, it fetches all numbers via GET /v1/direct_numbers/
// If numbers are provided, filters the results to only those numbers
func (c *Client) GetDirectNumbers(numbers ...string) ([]map[string]interface{}, error) {
	method := "/direct_numbers/"
	params := url.Values{}

	var resp struct {
		Status string                   `json:"status"`
		Info   []map[string]interface{} `json:"info"`
	}

	if err := c.Get(method, params, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	// If no specific numbers requested, return all
	if len(numbers) == 0 {
		return resp.Info, nil
	}

	// Filter to requested numbers only
	numberSet := make(map[string]bool)
	for _, num := range numbers {
		numberSet[num] = true
	}

	var result []map[string]interface{}
	for _, dn := range resp.Info {
		if num, ok := dn["number"].(string); ok && numberSet[num] {
			result = append(result, dn)
		}
	}

	// Verify all requested numbers were found
	if len(result) != len(numbers) {
		foundNumbers := make(map[string]bool)
		for _, dn := range result {
			if num, ok := dn["number"].(string); ok {
				foundNumbers[num] = true
			}
		}
		for _, num := range numbers {
			if !foundNumbers[num] {
				log.Debugf("Requested number not found: %s", num)
			}
		}
		return result, fmt.Errorf("not all requested numbers found")
	}

	return result, nil
}

// GetDirectCountries lists available direct number countries.
func (c *Client) GetDirectCountries() ([]map[string]interface{}, error) {
	method := "/direct_numbers/countries/"
	params := url.Values{}

	var resp struct {
		Status string                   `json:"status"`
		Info   []map[string]interface{} `json:"info"`
	}

	if err := c.Get(method, params, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	return resp.Info, nil
}

// GetDirectCountry lists destinations for a specific country.
func (c *Client) GetDirectCountry(country string) ([]map[string]interface{}, error) {
	method := "/direct_numbers/country/"
	params := url.Values{}
	params.Set("country", country)

	var resp struct {
		Status string                   `json:"status"`
		Info   []map[string]interface{} `json:"info"`
	}

	if err := c.Get(method, params, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	return resp.Info, nil
}

// GetDirectNumber returns information about a direct number.
func (c *Client) GetDirectNumber(number string) (map[string]interface{}, error) {
	// Validate phone number format
	if err := ValidatePhoneNumber(number); err != nil {
		return nil, err
	}

	method := "/direct_numbers/number/"
	params := url.Values{}
	params.Set("number", number)

	var resp struct {
		Status string                 `json:"status"`
		Data   map[string]interface{} `json:"data"`
		Info   map[string]interface{} `json:"info"`
	}

	if err := c.Get(method, params, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	// Support both response shapes: 'data' or 'info'
	if resp.Data != nil {
		return resp.Data, nil
	}
	if resp.Info != nil {
		return resp.Info, nil
	}

	return nil, fmt.Errorf("empty response")
}
