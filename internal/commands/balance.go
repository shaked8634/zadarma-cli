package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zadarma/zadarma-cli/internal/client"
)

// ClientFactory creates a new API client (injected by main).
type ClientFactory func() *client.Client

// NewBalanceCmd creates the 'balance' command.
func NewBalanceCmd(factory ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:          "balance",
		Short:        "Get account balance",
		Long:         "Get your Zadarma account balance",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			if !jsonOutput {
				of, _ := cmd.Root().PersistentFlags().GetString("output")
				jsonOutput = of == "json"
			}
			c := factory()

			balance, currency, err := c.GetBalance()
			if err != nil {
				return failCmd(cmd, err)
			}

			if jsonOutput {
				data := map[string]interface{}{
					"balance":  balance,
					"currency": currency,
				}
				out, _ := json.MarshalIndent(data, "", "  ")
				fmt.Println(string(out))
			} else {
				fmt.Printf("Balance: %.2f %s\n", balance, currency)
			}

			return nil
		},
	}
}
