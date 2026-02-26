package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewSMSCmd creates the 'sms' command group.
func NewSMSCmd(factory ClientFactory) *cobra.Command {
	var phoneNumber string
	var message string

	cmd := &cobra.Command{
		Use:   "sms",
		Short: "Send SMS messages",
		Long:  "Send SMS messages via Zadarma",
	}

	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send an SMS message",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")

			if phoneNumber == "" || message == "" {
				fmt.Fprintf(os.Stderr, "Error: --phone and --message are required\n")
				os.Exit(1)
			}

			c := factory()

			result, err := c.SendSMS(phoneNumber, message)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(out))
			} else {
				fmt.Printf("SMS Status: %s (Message ID: %s)\n", result["status"], result["id"])
			}
			return nil
		},
	}

	sendCmd.Flags().StringVar(&phoneNumber, "phone", "", "Recipient phone number")
	sendCmd.Flags().StringVar(&message, "message", "", "Message text")

	cmd.AddCommand(sendCmd)
	return cmd
}
