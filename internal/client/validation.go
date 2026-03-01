package client

import (
	"fmt"
	"net/url"
	"strconv"
)

// ValidateWebhookURL validates that the input is a valid webhook URL.
func ValidateWebhookURL(urlStr string) error {
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return fmt.Errorf("invalid webhook URL: %w", err)
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
