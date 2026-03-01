package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// NewWebhookCmd creates the 'webhook' command group.
func NewWebhookCmd(factory ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Manage Zadarma webhooks",
	}

	setCmd := &cobra.Command{
		Use:   "set [url]",
		Short: "Set the notification webhook URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			if !jsonOutput {
				of, _ := cmd.Root().PersistentFlags().GetString("output")
				jsonOutput = of == "json"
			}
			webhookURL := args[0]
			c := factory()

			result, err := c.SetWebhook(webhookURL)
			if err != nil {
				return err
			}

			if jsonOutput {
				b, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(b))
			} else {
				fmt.Printf("Webhook set successfully.\n")
				if v, ok := result["url"]; ok {
					fmt.Printf("URL: %v\n", v)
				}
			}
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get the current notification webhook URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			if !jsonOutput {
				of, _ := cmd.Root().PersistentFlags().GetString("output")
				jsonOutput = of == "json"
			}
			c := factory()

			result, err := c.GetWebhook()
			if err != nil {
				return err
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(out))
			} else {
				if v, ok := result["url"]; ok {
					fmt.Printf("Webhook URL: %v\n", v)
				} else {
					fmt.Println("No webhook URL configured.")
				}
			}
			return nil
		},
	}

	listenCmd := &cobra.Command{
		Use:   "listen",
		Short: "Listen for incoming SMS webhooks (Option A)",
		Long:  "Starts a local server to listen for and print incoming SMS from Zadarma webhooks.",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetString("port")
			if port == "" {
				port = "8080"
			}

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				// Handle Zadarma verification (zd_echo)
				if echo := r.URL.Query().Get("zd_echo"); echo != "" {
					_, _ = fmt.Fprint(w, echo)
					fmt.Printf("\n[VERIFICATION] Responded to zd_echo: %s\n", echo)
					return
				}

				if r.Method == "POST" {
					body, _ := io.ReadAll(r.Body)

					// Try to parse as SMS event
					var data map[string]interface{}
					if err := json.Unmarshal(body, &data); err == nil {
						if data["event"] == "SMS" {
							fmt.Printf("\n--- NEW SMS RECEIVED ---\n")
							fmt.Printf("From: %v\n", data["caller_id"])
							fmt.Printf("To:   %v\n", data["caller_did"])
							fmt.Printf("Text: %v\n", data["text"])
							fmt.Printf("------------------------\n")
						} else {
							fmt.Printf("\n[EVENT] %s: %s\n", data["event"], string(body))
						}
					} else {
						// Raw output if not JSON
						fmt.Printf("\n[WEBHOOK] Raw body: %s\n", string(body))
					}

					w.WriteHeader(http.StatusOK)
					_, _ = fmt.Fprint(w, "OK")
				}
			})

			fmt.Printf("🐧 Zadarma SMS Listener starting on port %s...\n", port)
			fmt.Println("Note: You must expose this port (e.g., via ngrok) and set the URL using 'zadarma webhook set'")
			return http.ListenAndServe(":"+port, nil)
		},
	}

	listenCmd.Flags().StringP("port", "p", "8080", "Port to listen on")

	cmd.AddCommand(setCmd, getCmd, listenCmd)
	return cmd
}
