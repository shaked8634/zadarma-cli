package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

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

	// Use a raw response container because API may return different shapes
	var raw map[string]json.RawMessage
	if err := c.Get(method, params, &raw); err != nil {
		return nil, err
	}

	// Check status field first if present
	if stRaw, ok := raw["status"]; ok {
		var st string
		_ = json.Unmarshal(stRaw, &st)
		if st != "success" {
			if st == "" {
				return nil, fmt.Errorf("API error: unknown status")
			}
			return nil, fmt.Errorf("API error: %s", st)
		}
	}

	// Normalize two possible shapes:
	// 1) {"data":[{"sender_id":"...","type":"..."}]}
	// 2) {"senders":["Teamsale", ...]}

	var result []map[string]interface{}

	if dataRaw, ok := raw["data"]; ok {
		// Try to unmarshal as []map[string]interface{}
		var arr []map[string]interface{}
		if err := json.Unmarshal(dataRaw, &arr); err == nil {
			result = arr
		} else {
			// If data present but not in expected shape, ignore and try senders
		}
	}

	if len(result) == 0 {
		if sendersRaw, ok := raw["senders"]; ok {
			var arr []string
			if err := json.Unmarshal(sendersRaw, &arr); err == nil {
				for _, s := range arr {
					result = append(result, map[string]interface{}{"sender_id": s, "type": ""})
				}
			}
		}
	}

	// If still empty, return error
	if len(result) == 0 {
		return nil, fmt.Errorf("no senders found or unexpected response shape")
	}

	return result, nil
}
