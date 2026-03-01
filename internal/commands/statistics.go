package commands

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// NewStatisticsCmd creates the 'statistics' command.
func NewStatisticsCmd(factory ClientFactory) *cobra.Command {
	var start, end string
	var sip int
	var costOnly bool

	cmd := &cobra.Command{
		Use:   "statistics",
		Short: "Get call statistics",
		Long:  "Get call statistics from your Zadarma account",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			const timeFormat = "2006-01-02 15:04:05"
			if start != "" {
				if _, err := time.Parse(timeFormat, start); err != nil {
					return fmt.Errorf("invalid start date format: %s. Expected YYYY-MM-DD HH:MM:SS", start)
				}
			}
			if end != "" {
				if _, err := time.Parse(timeFormat, end); err != nil {
					return fmt.Errorf("invalid end date format: %s. Expected YYYY-MM-DD HH:MM:SS", end)
				}
			}
			if cmd.Flags().Changed("sip") && sip <= 0 {
				return fmt.Errorf("invalid sip: %d. Must be a positive integer", sip)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")
			if !jsonOutput {
				if jb, _ := cmd.Root().PersistentFlags().GetBool("json"); jb {
					jsonOutput = true
				}
				of, _ := cmd.Root().PersistentFlags().GetString("output")
				jsonOutput = jsonOutput || of == "json"
			}
			c := factory()

			params := url.Values{}
			if start != "" {
				params.Set("start", start)
			}
			if end != "" {
				params.Set("end", end)
			}
			if sip != 0 {
				params.Set("sip", fmt.Sprintf("%d", sip))
			}
			if costOnly {
				params.Set("cost_only", "1")
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
	cmd.Flags().IntVar(&sip, "sip", 0, "SIP account to filter by (integer)")
	cmd.Flags().BoolVar(&costOnly, "cost-only", false, "Retrieve only costs")

	return cmd
}
