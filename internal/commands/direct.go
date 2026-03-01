package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// NewDirectCmd creates the 'direct' command group for virtual numbers.
func NewDirectCmd(factory ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "direct",
		Short: "Manage direct (virtual) numbers",
		Long:  "Commands for exploring direct number availability and details.",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "countries",
			Short: "List countries with direct numbers",
			RunE: func(cmd *cobra.Command, args []string) error {
				jsonOutput, _ := cmd.Flags().GetBool("json")
				if !jsonOutput {
					of, _ := cmd.Root().PersistentFlags().GetString("output")
					jsonOutput = of == "json"
				}
				c := factory()

				countries, err := c.GetDirectCountries()
				if err != nil {
					return failCmd(cmd, err)
				}

				if jsonOutput {
					out, _ := json.MarshalIndent(countries, "", "  ")
					fmt.Println(string(out))
				} else {
					for _, country := range countries {
						fmt.Printf("%v: %v\n", country["id"], country["name"])
					}
				}
				return nil
			},
		},
		&cobra.Command{
			Use:   "country <code>",
			Short: "List direct destinations for a specific country",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				jsonOutput, _ := cmd.Flags().GetBool("json")
				if !jsonOutput {
					of, _ := cmd.Root().PersistentFlags().GetString("output")
					jsonOutput = of == "json"
				}
				c := factory()

				countryData, err := c.GetDirectCountry(args[0])
				if err != nil {
					return failCmd(cmd, err)
				}

				if jsonOutput {
					out, _ := json.MarshalIndent(countryData, "", "  ")
					fmt.Println(string(out))
				} else {
					if len(countryData) == 0 {
						fmt.Printf("No direct numbers found for %s\n", args[0])
						return nil
					}
					for _, entry := range countryData {
						fmt.Printf("%v: %v\n", entry["id"], entry["name"])
					}
				}
				return nil
			},
		},
		&cobra.Command{
			Use:   "number <type> <number>",
			Short: "Get information for a specific direct number",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				jsonOutput, _ := cmd.Flags().GetBool("json")
				if !jsonOutput {
					of, _ := cmd.Root().PersistentFlags().GetString("output")
					jsonOutput = of == "json"
				}
				c := factory()

				info, err := c.GetDirectNumber(args[0], args[1])
				if err != nil {
					return failCmd(cmd, err)
				}

				if jsonOutput {
					out, _ := json.MarshalIndent(info, "", "  ")
					fmt.Println(string(out))
				} else {
					for key, value := range info {
						fmt.Printf("%s: %v\n", key, value)
					}
				}
				return nil
			},
		},
	)

	return cmd
}
