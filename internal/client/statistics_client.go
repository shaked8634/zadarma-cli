package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// StatisticResponse represents the data returned by the statistics API.
type StatisticResponse struct {
	Status string                   `json:"status"`
	Data   []map[string]interface{} `json:"data"`
	Stats  []map[string]interface{} `json:"stats"`
	Start  string                   `json:"start"`
	End    string                   `json:"end"`
}

// GetStatistics fetches call statistics.
// Possible parameters: start, end, user_id, sip, etc.
func (c *Client) GetStatistics(params url.Values) ([]map[string]interface{}, error) {
	method := "/statistics/"

	var raw map[string]interface{}

	if err := c.Get(method, params, &raw); err != nil {
		return nil, err
	}

	if st, _ := raw["status"].(string); st != "success" {
		return nil, fmt.Errorf("API error: %s", st)
	}

	// Try "data" field first, then "stats"
	var stats []map[string]interface{}

	// Handle "data" field - may be []interface{} or []map[string]interface{}
	if dataRaw, ok := raw["data"].([]interface{}); ok {
		for _, item := range dataRaw {
			if m, ok := item.(map[string]interface{}); ok {
				stats = append(stats, m)
			}
		}
	} else if data, ok := raw["data"].([]map[string]interface{}); ok {
		stats = data
	}

	// If no data, try stats field
	if len(stats) == 0 {
		if statsRaw, ok := raw["stats"].([]interface{}); ok {
			for _, s := range statsRaw {
				if m, ok := s.(map[string]interface{}); ok {
					stats = append(stats, m)
				}
			}
		}
	}

	// Normalize JSON types
	b, _ := json.Marshal(stats)
	var norm []map[string]interface{}
	_ = json.Unmarshal(b, &norm)

	return norm, nil
}
