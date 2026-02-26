package client

import (
	"fmt"
	"net/url"
)

// StatisticResponse represents the data returned by the statistics API.
type StatisticResponse struct {
	Status string                   `json:"status"`
	Data   []map[string]interface{} `json:"data"`
}

// GetStatistics fetches call statistics.
// Possible parameters: start, end, user_id, sip, etc.
func (c *Client) GetStatistics(params url.Values) ([]map[string]interface{}, error) {
	method := "/statistics/"

	var resp StatisticResponse

	if err := c.Get(method, params, &resp); err != nil {
		return nil, err
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	return resp.Data, nil
}
