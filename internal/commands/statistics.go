package commands

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

// NewStatisticsCmd creates the 'statistics' command.
func NewStatisticsCmd(factory ClientFactory) *cobra.Command {
	var start, end, sip string

	cmd := &cobra.Command{
		Use:   "statistics",
		Short: "Get call statistics",
		Long:  "Get call statistics from your Zadarma account",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			c := factory()

			params := url.Values{}
			if start != "" {
				params.Set("start", start)
			}
			if end != "" {
				params.Set("end", end)
			}
			if sip != "" {
				params.Set("sip", sip)
			}

			stats, err := c.GetStatistics(params)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(stats, "", "  ")
				fmt.Println(string(out))
			} else {
				if len(stats) == 0 {
					fmt.Println("No statistics found for the given period.")
					return nil
				}

				fmt.Printf("%-20s | %-15s | %-10s | %-10s | %-s\n", "Time", "From", "Duration", "Cost", "Destination")
				fmt.Println("--------------------------------------------------------------------------------")
				for _, s := range stats {
					fmt.Printf("%-20v | %-15v | %-10v | %-10v | %-v\n",
						s["callstart"],
						s["sip"],
						s["billseconds"],
						s["billcost"],
						s["destination"],
					)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&start, "start", "", "Start date (YYYY-MM-DD HH:MM:SS)")
	cmd.Flags().StringVar(&end, "end", "", "End date (YYYY-MM-DD HH:MM:SS)")
	cmd.Flags().StringVar(&sip, "sip", "", "SIP account to filter by")

	return cmd
}
