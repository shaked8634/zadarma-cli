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
	"github.com/zadarma/zadarma-cli/internal/log"
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
Before starting, fetches the current webhook URL and displays it.
Listens on specified port and prints received SMS in text or JSON format.
Can be backgrounded with '&' to keep running while using other CLI commands.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			port, _ := cmd.Flags().GetString("port")
			if port == "" {
				port = "8080"
			}

			// Optional: set webhook URL before starting
			setURL, _ := cmd.Flags().GetString("set-webhook")
			enableSMS, _ := cmd.Flags().GetBool("enable-sms")

			// Fetch current webhook URL first
			c := factory()
			webhookInfo, err := c.GetWebhooks()
			if err != nil {
				log.Debugf("Failed to fetch webhook info: %v", err)
				_, _ = fmt.Fprintf(os.Stderr, "Warning: Could not fetch webhook URL: %v\n", err)
			} else {
				if currentURL, ok := webhookInfo["url"]; ok && currentURL != "" {
					fmt.Printf("Current webhook URL: %v\n", currentURL)
				} else {
					fmt.Println("No webhook URL currently configured.")
				}
			}

			// If user requested to set a webhook URL, do it now
			if setURL != "" {
				fmt.Printf("Setting webhook URL to %s...\n", setURL)
				res, err := c.SetWebhook(setURL)
				if err != nil {
					return failCmd(cmd, fmt.Errorf("failed to set webhook URL: %w", err))
				}
				if jsonOutput {
					out, _ := json.MarshalIndent(res, "", "  ")
					fmt.Println(string(out))
				} else {
					fmt.Printf("Webhook set: %v\n", res)
				}
			}

			// If user wants to enable SMS hooks, call the API
			if enableSMS {
				fmt.Println("Enabling SMS webhook hooks...")
				res, err := c.SetWebhookHooks(true)
				if err != nil {
					return failCmd(cmd, fmt.Errorf("failed to enable sms hooks: %w", err))
				}
				if jsonOutput {
					out, _ := json.MarshalIndent(res, "", "  ")
					fmt.Println(string(out))
				} else {
					fmt.Printf("Webhook hooks updated: %v\n", res)
				}
			}

			fmt.Printf("Listening for SMS webhooks on port %s...\n", port)
			fmt.Println("Press Ctrl+C to stop.")
			fmt.Println("Note: You must expose this port via ngrok or similar and set it with 'sms webhook set <url>'")
			fmt.Println()

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
		},
	}

	listenCmd.Flags().StringP("port", "p", "8080", "Port to listen on")
	listenCmd.Flags().String("set-webhook", "", "Set webhook URL before starting listener")
	listenCmd.Flags().Bool("enable-sms", false, "Enable SMS webhook event hooks before starting listener")

	cmd.AddCommand(sendCmd)
	cmd.AddCommand(sendersCmd)
	cmd.AddCommand(listenCmd)
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
