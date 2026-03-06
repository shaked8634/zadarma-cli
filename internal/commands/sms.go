package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

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
		Short: "Send SMS messages and listen for incoming webhooks",
		Long:  "Send SMS messages via Zadarma and receive incoming SMS via webhooks",
		// Do not show Cobra usage when runtime/API errors occur; usage should be shown only for CLI syntax errors
		SilenceUsage: true,
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
		Use:   "senders [phones]",
		Short: "Get valid SMS senders for given phone numbers",
		Long:  "Get valid SMS senders for a comma-separated list of phone numbers or pass them as positional arguments",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			if !jsonOutput {
				if jb, _ := cmd.Root().PersistentFlags().GetBool("json"); jb {
					jsonOutput = true
				}
				of, _ := cmd.Root().PersistentFlags().GetString("output")
				jsonOutput = jsonOutput || of == "json"
			}

			// If positional args provided, join them into phones
			if len(args) > 0 {
				phones = strings.Join(args, ",")
			}

			if phones == "" {
				return failCmd(cmd, fmt.Errorf("--phones is required (or pass phone numbers as args)"))
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
				// Determine if any sender has a non-empty type
				hasType := false
				for _, s := range result {
					if t, ok := s["type"]; ok && t != nil && fmt.Sprint(t) != "" {
						hasType = true
						break
					}
				}
				if hasType {
					fmt.Printf("%-20s | %-15s\n", "Sender", "Type")
					fmt.Println("---------------------------------------")
					for _, s := range result {
						fmt.Printf("%-20v | %-15v\n", s["sender_id"], s["type"])
					}
				} else {
					// Single column list
					for _, s := range result {
						fmt.Println(s["sender_id"])
					}
				}
			}
			return nil
		},
	}
	sendersCmd.Flags().StringVar(&phones, "phones", "", "Destination phone numbers (comma separated)")

	listenCmd := &cobra.Command{
		Use:   "listen",
		Short: "Listen for incoming SMS webhooks (daemon mode)",
		Long: `Start a local HTTP server to receive incoming SMS webhooks from Zadarma.

Before starting, a webhook URL must already be configured in Zadarma.
Listens on specified port and prints received SMS in text or JSON format.
Can be run in the background with '&' to keep running while using other CLI commands.

To configure a webhook URL first, use:
  zadarma-cli sms set-webhook <URL>

Examples:
  # Listen on default port 8080
  zadarma-cli sms listen
  
  # Listen on custom port
  zadarma-cli sms listen --port 9000`,
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			port, _ := cmd.Flags().GetString("port")
			if port == "" {
				port = "8080"
			}

			// New: webhook flag to register a webhook before listening
			webhookFlag, _ := cmd.Flags().GetString("webhook")

			c := factory()

			if webhookFlag != "" {
				// Register the webhook and enable SMS hooks
				fmt.Printf("Registering webhook URL: %s\n", webhookFlag)
				if _, err := c.SetWebhook(webhookFlag); err != nil {
					return failCmd(cmd, fmt.Errorf("failed to set webhook URL: %w", err))
				}
				fmt.Println("✓ Webhook registered")
				if _, err := c.SetWebhookHooks(true); err != nil {
					return failCmd(cmd, fmt.Errorf("failed to enable sms hooks: %w", err))
				}
				fmt.Println("✓ SMS webhooks enabled")
			}

			// Fetch current webhook URL first
			webhookInfo, err := c.GetWebhooks()
			if err != nil {
				return failCmd(cmd, fmt.Errorf("failed to fetch webhook info: %w", err))
			}

			// Check if webhook is configured
			currentURL, ok := webhookInfo["url"]
			if !ok || currentURL == nil || fmt.Sprint(currentURL) == "" {
				return failCmd(cmd, fmt.Errorf("no webhook URL configured in Zadarma. Use 'zadarma-cli sms set-webhook <URL>' to configure one"))
			}

			fmt.Printf("Using webhook URL: %v\n", currentURL)
			fmt.Printf("\n🐧 Zadarma SMS Listener starting on port %s...\n", port)
			fmt.Println("Waiting for incoming SMS... (Press Ctrl+C to stop)")
			fmt.Println()

			return startSMSListener(port, jsonOutput)
		},
	}

	listenCmd.Flags().StringP("port", "p", "8080", "Port to listen on")
	listenCmd.Flags().StringP("webhook", "w", "", "Webhook URL to register before listening")
	// Silence usage for listen command
	listenCmd.SilenceUsage = true

	setWebhookCmd := &cobra.Command{
		Use:   "set-webhook <WEBHOOK>",
		Short: "Set webhook URL and start listening for SMS",
		Long: `Register a webhook URL with Zadarma and immediately start listening for incoming SMS.

This command:
  1. Validates the webhook URL
  2. Registers it with Zadarma
  3. Enables SMS webhook notifications
  4. Starts listening on the specified port

The listener runs in the foreground. You can run other CLI commands in separate terminals.

Examples:
  # Using localtunnel
  zadarma-cli sms set-webhook https://my-tunnel.loca.lt
  
  # Using ngrok
  zadarma-cli sms set-webhook https://abc123.ngrok.io
  
  # Custom port
  zadarma-cli sms set-webhook https://my-tunnel.loca.lt --port 9000`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			webhookURL := args[0]
			port, _ := cmd.Flags().GetString("port")
			if port == "" {
				port = "8080"
			}

			c := factory()

			// Register the webhook URL
			fmt.Printf("Registering webhook URL: %s\n", webhookURL)
			res, err := c.SetWebhook(webhookURL)
			if err != nil {
				return failCmd(cmd, fmt.Errorf("failed to set webhook URL: %w", err))
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(res, "", "  ")
				fmt.Println(string(out))
			} else {
				fmt.Printf("✓ Webhook registered\n")
			}

			// Enable SMS webhook notifications
			fmt.Println("Enabling SMS webhook notifications...")
			res, err = c.SetWebhookHooks(true)
			if err != nil {
				return failCmd(cmd, fmt.Errorf("failed to enable sms hooks: %w", err))
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(res, "", "  ")
				fmt.Println(string(out))
			} else {
				fmt.Printf("✓ SMS webhooks enabled\n")
			}

			fmt.Printf("\n🐧 Zadarma SMS Listener starting on port %s...\n", port)
			fmt.Println("Waiting for incoming SMS... (Press Ctrl+C to stop)")
			fmt.Println()

			return startSMSListener(port, jsonOutput)
		},
	}

	setWebhookCmd.Flags().StringP("port", "p", "8080", "Port to listen on")
	// Silence usage for set-webhook command
	setWebhookCmd.SilenceUsage = true

	cmd.AddCommand(sendCmd)
	cmd.AddCommand(sendersCmd)
	cmd.AddCommand(listenCmd)
	cmd.AddCommand(setWebhookCmd)
	return cmd
}

// printSMSEvent formats an incoming SMS event for text output
func printSMSEvent(data map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() { _ = w.Flush() }()

	_, _ = fmt.Fprintln(w, "--- INCOMING SMS ---")
	_, _ = fmt.Fprintln(w, "FROM\t", getStringField(data, "caller_id", "-"))
	_, _ = fmt.Fprintln(w, "TO\t", getStringField(data, "caller_did", "-"))
	_, _ = fmt.Fprintln(w, "TEXT\t", getStringField(data, "text", "-"))
	if ts, ok := data["timestamp"]; ok {
		_, _ = fmt.Fprintln(w, "TIME\t", fmt.Sprint(ts))
	}
	_, _ = fmt.Fprintln(w, "-------------------")
}

// startSMSListener starts an HTTP server to listen for incoming SMS webhooks
func startSMSListener(port string, jsonOutput bool) error {
	// HTTP handler for incoming SMS
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle Zadarma verification (zd_echo)
		if echo := r.URL.Query().Get("zd_echo"); echo != "" {
			_, _ = fmt.Fprint(w, echo)
			if !jsonOutput {
				fmt.Printf("[VERIFICATION] Responded to zd_echo: %s\n", echo)
			}
			return
		}

		if r.Method == "POST" {
			body, _ := io.ReadAll(r.Body)
			_ = r.Body.Close()

			var data map[string]interface{}
			parseErr := json.Unmarshal(body, &data)

			if parseErr == nil && data["event"] == "SMS" {
				// SMS event received
				if jsonOutput {
					// JSON output: just print the parsed data
					out, _ := json.MarshalIndent(data, "", "  ")
					fmt.Println(string(out))
				} else {
					// Text output: formatted table
					printSMSEvent(data)
				}
			} else if parseErr == nil {
				// Other event types
				if jsonOutput {
					out, _ := json.MarshalIndent(data, "", "  ")
					fmt.Println(string(out))
				} else {
					fmt.Printf("[EVENT] Type: %v\n", data["event"])
				}
			} else {
				// Raw output if not JSON-parseable
				if jsonOutput {
					// In JSON mode, print what we got as a structured error message
					errData := map[string]interface{}{
						"error": "Failed to parse webhook body as JSON",
						"raw":   string(body),
					}
					out, _ := json.MarshalIndent(errData, "", "  ")
					fmt.Println(string(out))
				} else {
					fmt.Printf("[WEBHOOK] Raw body: %s\n", string(body))
				}
			}

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, "OK")
		}
	})

	return http.ListenAndServe(":"+port, nil)
}
