package client

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ValidateWebhookURL validates that the input is a valid webhook URL.
func ValidateWebhookURL(urlStr string) error {
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return fmt.Errorf("invalid webhook URL: %w", err)
	}

	return nil
}

// ValidateWebhookEcho sends a verification request to the webhook URL and checks
// that it returns the expected echo value. This is required by Zadarma to verify
// ownership of the URL before accepting it as a webhook endpoint.
func ValidateWebhookEcho(webhookURL string) error {
	randCode := fmt.Sprintf("%d", rand.Intn(10000000)+1000000)

	testURL, err := url.Parse(webhookURL)
	if err != nil {
		return fmt.Errorf("failed to parse webhook URL: %w", err)
	}

	q := testURL.Query()
	q.Set("zd_echo", randCode)
	testURL.RawQuery = q.Encode()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(testURL.String())
	if err != nil {
		return fmt.Errorf("failed to reach webhook URL: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	returned := strings.TrimSpace(string(body))
	if returned != randCode {
		return fmt.Errorf("validation failed: expected '%s', got '%s'", randCode, returned)
	}

	return nil
}

// ValidatePhoneNumber validates that the input is a valid phone number.
// It must be a valid uint and be between 8 and 20 digits (inclusive).
func ValidatePhoneNumber(phoneNumber string) error {
	if phoneNumber == "" {
		return fmt.Errorf("phone number cannot be empty")
	}

	// Parse as uint to ensure it's a valid number (all digits, no special chars)
	_, err := strconv.ParseUint(phoneNumber, 10, 64)
	if err != nil {
		return fmt.Errorf("phone number must contain only digits, got: %s", phoneNumber)
	}

	// Check length
	length := len(phoneNumber)
	if length < 8 {
		return fmt.Errorf("phone number must be at least 8 digits, got %d digits", length)
	}
	if length > 20 {
		return fmt.Errorf("phone number must be at most 20 digits, got %d digits", length)
	}

	return nil
}
