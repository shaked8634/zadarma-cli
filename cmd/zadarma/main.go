package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zadarma/zadarma-cli/internal/client"
	"github.com/zadarma/zadarma-cli/internal/commands"
)

const Version = "0.0.5"

var (
	apiKey       string
	apiSecret    string
	outputFormat string
	debug        bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "zadarma-cli",
		Short:   "Zadarma VoIP API command-line client",
		Long:    "zadarma-cli - A CLI tool for interacting with the Zadarma VoIP API",
		Version: Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if apiKey == "" {
				apiKey = os.Getenv("ZADARMA_API_KEY")
			}
			if apiSecret == "" {
				apiSecret = os.Getenv("ZADARMA_API_SECRET")
			}

			if apiKey == "" || apiSecret == "" {
				_, _ = fmt.Fprintf(os.Stderr, "Error: ZADARMA_API_KEY and ZADARMA_API_SECRET must be set\n")
				_, _ = fmt.Fprintf(os.Stderr, "Use --key and --secret flags or set environment variables\n")
				os.Exit(1)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&apiKey, "key", "k", "", "Zadarma API key")
	rootCmd.PersistentFlags().StringVarP(&apiSecret, "secret", "s", "", "Zadarma API secret")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")

	// Hide the explicit 'help' subcommand (users can still use --help)
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	clientFactory := func() *client.Client {
		return client.NewClient(apiKey, apiSecret, debug)
	}

	rootCmd.AddCommand(commands.NewBalanceCmd(clientFactory))
	rootCmd.AddCommand(commands.NewSIPCmd(clientFactory))
	rootCmd.AddCommand(commands.NewPhoneCmd(clientFactory))
	rootCmd.AddCommand(commands.NewSMSCmd(clientFactory))
	rootCmd.AddCommand(commands.NewPBXCmd(clientFactory))
	rootCmd.AddCommand(commands.NewStatisticsCmd(clientFactory))
	rootCmd.AddCommand(commands.NewWebhookCmd(clientFactory))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
