package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// IncomingSMS represents a received SMS from Zadarma webhook
type IncomingSMS struct {
	Event     string `json:"event"`
	CallerID  string `json:"caller_id"`
	CallerDid string `json:"caller_did"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp,omitempty"`
}

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
				// Show useful fields from the full response
				if messages, ok := result["messages"].(float64); ok {
					fmt.Printf("Messages sent: %.0f\n", messages)
				}
				if cost, ok := result["cost"].(float64); ok {
					currency, _ := result["currency"].(string)
					fmt.Printf("Cost: %.2f %s\n", cost, currency)
				}
				if det, ok := result["sms_detalization"].([]interface{}); ok && len(det) > 0 {
					if det0, ok := det[0].(map[string]interface{}); ok {
						if parts, ok := det0["parts"].(float64); ok {
							fmt.Printf("Parts: %.0f\n", parts)
						}
						if msg, ok := det0["message"].(string); ok {
							fmt.Printf("Message: %s\n", msg)
						}
					}
				}
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
		Short: "Listen for incoming SMS webhooks",
		Long: `Start a local HTTP server to receive incoming SMS webhooks from Zadarma.

If --webhook is provided:
  1. Start local HTTP server
  2. Register the webhook URL with Zadarma
  3. Wait for Zadarma validation (zd_echo)
  4. Continue listening for SMS

If no --webhook is provided:
  1. Check if a webhook is already configured
  2. If yes, start listening for SMS
  3. If no, exit with error

Examples:
  # With a webhook URL (requires tunnel like ngrok/localtunnel)
  zadarma-cli sms listen --webhook https://abc123.ngrok.io
  
  # Using existing webhook
  zadarma-cli sms listen
  
  # Custom port
  zadarma-cli sms listen --port 9000`,
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			port, _ := cmd.Flags().GetString("port")
			if port == "" {
				port = "8080"
			}

			webhookFlag, _ := cmd.Flags().GetString("webhook")

			c := factory()

			// Start HTTP server for listening
			validateChan := make(chan string, 1)
			serverErr := make(chan error, 1)

			go func() {
				serverErr <- startSMSListener(port, jsonOutput, validateChan)
			}()

			// Give server time to start
			time.Sleep(500 * time.Millisecond)

			if webhookFlag != "" {
				// Case A: Register new webhook and wait for validation
				fmt.Printf("Registering webhook URL: %s\n", webhookFlag)
				if _, err := c.SetWebhook(webhookFlag); err != nil {
					return failCmd(cmd, fmt.Errorf("failed to set webhook URL: %w", err))
				}
				fmt.Println("✓ Webhook registered")

				// Check current hooks state
				webhookInfo, _ := c.GetWebhooks()
				hooks, _ := webhookInfo["hooks"].(map[string]interface{})
				smsEnabled := false
				if hooks != nil {
					if v, ok := hooks["sms"]; ok {
						smsEnabled = v == "true" || v == true
					}
				}

				if !smsEnabled {
					if _, err := c.SetWebhookHooks(true); err != nil {
						errStr := err.Error()
						if !strings.Contains(errStr, "Update error") {
							return failCmd(cmd, fmt.Errorf("failed to enable sms hooks: %w", err))
						}
					}
				}
				fmt.Println("✓ SMS webhooks enabled")

				// Wait for validation (zd_echo) with 60 second timeout
				fmt.Println("[INFO] Waiting for Zadarma validation (zd_echo request)...")
				fmt.Println("[INFO] Timeout: 60 seconds")

				select {
				case echoValue := <-validateChan:
					fmt.Printf("Received zd_echo: %s\n", echoValue)
					fmt.Println("Webhook validated successfully!")
				case <-time.After(60 * time.Second):
					return failCmd(cmd, fmt.Errorf("validation timeout: Zadarma did not send zd_echo within 60 seconds. Check that your URL is accessible"))
				}
			} else {
				// Case B: Use existing webhook
				webhookInfo, err := c.GetWebhooks()
				if err != nil {
					return failCmd(cmd, fmt.Errorf("failed to fetch webhook info: %w", err))
				}

				currentURL, ok := webhookInfo["url"]
				if !ok || currentURL == nil || fmt.Sprint(currentURL) == "" {
					return failCmd(cmd, fmt.Errorf("no webhook URL configured. Run 'zadarma sms set-webhook <URL>' first to configure a webhook"))
				}

				fmt.Printf("Using webhook URL: %v\n", currentURL)
			}

			fmt.Printf("\n🐧 Zadarma SMS Listener starting on port %s...\n", port)
			fmt.Println("Waiting for incoming SMS... (Press Ctrl+C to stop)")
			fmt.Println()

			// Return error from server if any
			return <-serverErr
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

			return startSMSListener(port, jsonOutput, nil)
		},
	}

	setWebhookCmd.Flags().StringP("port", "p", "8080", "Port to listen on")
	// Silence usage for set-webhook command
	setWebhookCmd.SilenceUsage = true

	getWebhookCmd := &cobra.Command{
		Use:   "get-webhook",
		Short: "Get the current notification webhook URL",
		Long:  "Retrieves the currently configured webhook URL from Zadarma.",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			c := factory()

			result, err := c.GetWebhooks()
			if err != nil {
				return failCmd(cmd, err)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(out))
			} else {
				if v, ok := result["url"]; ok && v != nil && v != "" {
					fmt.Printf("Webhook URL: %v\n", v)
					if hooks, ok := result["hooks"].(map[string]interface{}); ok {
						fmt.Println("Enabled hooks:")
						for hook, enabled := range hooks {
							fmt.Printf("  - %s: %v\n", hook, enabled)
						}
					}
				} else {
					fmt.Println("No webhook URL configured.")
				}
			}
			return nil
		},
		SilenceUsage: true,
	}

	cmd.AddCommand(sendCmd)
	cmd.AddCommand(sendersCmd)
	cmd.AddCommand(listenCmd)
	cmd.AddCommand(setWebhookCmd)
	cmd.AddCommand(getWebhookCmd)
	return cmd
}

// printSMSEvent formats an incoming SMS event for text output
func printSMSEvent(sms *IncomingSMS) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() { _ = w.Flush() }()

	_, _ = fmt.Fprintln(w, "--- INCOMING SMS ---")
	_, _ = fmt.Fprintln(w, "FROM\t", sms.CallerID)
	_, _ = fmt.Fprintln(w, "TO\t", sms.CallerDid)
	_, _ = fmt.Fprintln(w, "TEXT\t", sms.Text)
	if sms.Timestamp != "" {
		_, _ = fmt.Fprintln(w, "TIME\t", sms.Timestamp)
	}
	_, _ = fmt.Fprintln(w, "-------------------")
}

// startSMSListener starts an HTTP server to listen for incoming SMS webhooks.
// If validateChan is not nil, it will be sent the zd_echo value when validation request is received.
func startSMSListener(port string, jsonOutput bool, validateChan chan string) error {
	// HTTP handler for incoming SMS
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle Zadarma verification (zd_echo)
		if echo := r.URL.Query().Get("zd_echo"); echo != "" {
			_, _ = fmt.Fprint(w, echo)

			// Send to channel if provided
			if validateChan != nil {
				select {
				case validateChan <- echo:
				default:
				}
			}

			fmt.Printf("[DEBUG] zd_echo request received: %s\n", echo)
			fmt.Printf("[DEBUG] Responding with: %s\n", echo)
			fmt.Printf("[DEBUG] Full response headers: Content-Type=%s\n", w.Header().Get("Content-Type"))
			fmt.Println("[INFO] Webhook validation successful! Continuing to listen for SMS...")
			return
		}

		if r.Method == "POST" {
			body, _ := io.ReadAll(r.Body)
			_ = r.Body.Close()

			// Parse form-encoded body (e.g., event=SMS&result=%7B...%7D)
			form, parseErr := url.ParseQuery(string(body))
			if parseErr != nil {
				fmt.Printf("[ERROR] Failed to parse form data: %v\n", parseErr)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			event := form.Get("event")

			if event == "SMS" {
				// Get the URL-encoded JSON result
				resultEncoded := form.Get("result")
				if resultEncoded == "" {
					fmt.Println("[ERROR] Missing 'result' field in SMS event")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				// Decode URL-encoded JSON
				resultJSON, decodeErr := url.QueryUnescape(resultEncoded)
				if decodeErr != nil {
					fmt.Printf("[ERROR] Failed to decode result: %v\n", decodeErr)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				// Parse inner JSON
				var sms IncomingSMS
				if jsonErr := json.Unmarshal([]byte(resultJSON), &sms); jsonErr != nil {
					fmt.Printf("[ERROR] Failed to parse SMS JSON: %v\n", jsonErr)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if jsonOutput {
					out, _ := json.MarshalIndent(sms, "", "  ")
					fmt.Println(string(out))
				} else {
					printSMSEvent(&sms)
				}
			} else if event != "" {
				// Other event types
				if jsonOutput {
					eventData := map[string]string{"event": event}
					out, _ := json.MarshalIndent(eventData, "", "  ")
					fmt.Println(string(out))
				} else {
					fmt.Printf("[EVENT] %s\n", event)
				}
			} else {
				// No event field - raw output
				if jsonOutput {
					errData := map[string]interface{}{
						"error": "No 'event' field in webhook body",
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
