package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// NewSMSCmd creates the 'sms' command group.
func NewSMSCmd(factory ClientFactory) *cobra.Command {
	var phoneNumber string
	var message string
	var sender string
	var phones string

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
			if !jsonOutput {
				if jb, _ := cmd.Root().PersistentFlags().GetBool("json"); jb {
					jsonOutput = true
				}
				of, _ := cmd.Root().PersistentFlags().GetString("output")
				jsonOutput = jsonOutput || of == "json"
			}

			if phoneNumber == "" || message == "" {
				return failCmd(cmd, fmt.Errorf("--phone and --message are required"))
			}

			c := factory()

			result, err := c.SendSMS(phoneNumber, message, sender)
			if err != nil {
				return failCmd(cmd, err)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(out))
			} else {
				// Safely extract normalized fields
				status := ""
				if v, ok := result["status"]; ok && v != nil {
					status = fmt.Sprint(v)
				} else {
					// If missing, assume success because API returned HTTP 200 and client validated status
					status = "success"
				}
				msgID := ""
				if v, ok := result["id"]; ok && v != nil {
					msgID = fmt.Sprint(v)
				}
				fmt.Printf("SMS Status: %s", status)
				if msgID != "" {
					fmt.Printf(" (Message ID: %s)", msgID)
				}
				fmt.Println()
			}
			return nil
		},
	}

	sendCmd.Flags().StringVar(&phoneNumber, "phone", "", "Recipient phone number")
	sendCmd.Flags().StringVar(&message, "message", "", "Message text")
	sendCmd.Flags().StringVar(&sender, "sender", "", "SMS sender (virtual number or text)")

	sendersCmd := &cobra.Command{
		Use:   "senders",
		Short: "Get valid SMS senders for a given phone number",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			if !jsonOutput {
				if jb, _ := cmd.Root().PersistentFlags().GetBool("json"); jb {
					jsonOutput = true
				}
				of, _ := cmd.Root().PersistentFlags().GetString("output")
				jsonOutput = jsonOutput || of == "json"
			}

			if phones == "" {
				return failCmd(cmd, fmt.Errorf("--phones is required"))
			}

			c := factory()

			result, err := c.GetSMSSenders(phones)
			if err != nil {
				return failCmd(cmd, err)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(out))
			} else {
				if len(result) == 0 {
					fmt.Println("No senders found.")
					return nil
				}
				fmt.Printf("%-20s | %-15s\n", "Sender", "Type")
				fmt.Println("---------------------------------------")
				for _, s := range result {
					fmt.Printf("%-20v | %-15v\n", s["sender_id"], s["type"])
				}
			}
			return nil
		},
	}
	sendersCmd.Flags().StringVar(&phones, "phones", "", "Destination phone numbers (comma separated)")

	cmd.AddCommand(sendCmd)
	cmd.AddCommand(sendersCmd)
	return cmd
}
