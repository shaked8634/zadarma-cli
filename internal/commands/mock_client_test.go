package commands

// MockClient is a mock client for testing command handlers.
type MockClient struct {
	balance  interface{}
	currency string
	err      error
}

func (m *MockClient) GetBalance() (interface{}, string, error) {
	return m.balance, m.currency, m.err
}

func (m *MockClient) GetPrice(number string) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return map[string]interface{}{"price": "0.05", "currency": "USD"}, nil
}

func (m *MockClient) GetSIPs() ([]map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []map[string]interface{}{
		{"sip_user": "user1", "status": "active"},
	}, nil
}

func (m *MockClient) GetSIPStatus(id string) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return true, nil
}

func (m *MockClient) SendSMS(phoneNumber, message, sender string) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return map[string]interface{}{"id": "msg123", "status": "sent"}, nil
}

func (m *MockClient) GetSMSSenders(phones string) ([]map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []map[string]interface{}{
		{"sender_id": "Zadarma", "type": "alpha"},
	}, nil
}

func (m *MockClient) GetDirectNumbers(numbers ...string) ([]map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []map[string]interface{}{
		{"number": "972556620707", "country": "Israel", "status": "on"},
	}, nil
}

func (m *MockClient) GetDirectCountries() ([]map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []map[string]interface{}{
		{"id": "US", "name": "United States"},
	}, nil
}

func (m *MockClient) GetDirectCountry(country string) ([]map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []map[string]interface{}{
		{"id": "212", "name": "New York"},
	}, nil
}

func (m *MockClient) GetDirectNumber(number string) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return map[string]interface{}{"number": number, "status": "available"}, nil
}

func (m *MockClient) GetPBXInfo(pbxID, numbers string) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return map[string]interface{}{"name": "My PBX", "status": "active"}, nil
}

func (m *MockClient) GetPBXInternalStatus(pbxID string) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return map[string]interface{}{"pbx_id": pbxID, "is_online": "true"}, nil
}

func (m *MockClient) GetPBXInternalInfo(pbxID string) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return map[string]interface{}{"pbx_id": pbxID, "title": "Main Office"}, nil
}

func (m *MockClient) SetWebhook(urlStr string) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return map[string]interface{}{"url": urlStr, "status": "set"}, nil
}

func (m *MockClient) GetWebhook() (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return map[string]interface{}{"url": "https://example.com/webhook", "status": "active"}, nil
}

func (m *MockClient) GetStatistics(startTime, endTime, sip string, costOnly bool) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return map[string]interface{}{"total_cost": "10.50", "calls": 100}, nil
}
