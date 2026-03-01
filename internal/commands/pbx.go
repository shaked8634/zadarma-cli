package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// NewPBXCmd creates the 'pbx' command group.
func NewPBXCmd(factory ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pbx",
		Short: "Manage PBX configuration",
		Long:  "PBX configuration and management commands",
	}

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Get PBX information",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			c := factory()

			pbxID, _ := cmd.Flags().GetString("pbx-id")
			numbers, _ := cmd.Flags().GetString("numbers")

			pbxInfo, err := c.GetPBXInfo(pbxID, numbers)
			if err != nil {
				return failCmd(cmd, err)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(pbxInfo, "", "  ")
				fmt.Println(string(out))
				return nil
			}

			// Text output: print a simple table with pbx_id and numbers
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			defer func() { _ = w.Flush() }()

			fmt.Fprintln(w, "pbx_id\tnumbers")
			fmt.Fprintln(w, "------\t-------")

			pbxIDVal := ""
			if v, ok := pbxInfo["pbx_id"]; ok {
				pbxIDVal = fmt.Sprintf("%v", v)
			}

			numbersVal := ""
			if v, ok := pbxInfo["numbers"]; ok {
				switch t := v.(type) {
				case []interface{}:
					parts := make([]string, 0, len(t))
					for _, it := range t {
						parts = append(parts, fmt.Sprintf("%v", it))
					}
					numbersVal = strings.Join(parts, ",")
				default:
					numbersVal = fmt.Sprintf("%v", v)
				}
			}

			fmt.Fprintf(w, "%s\t%s\n", pbxIDVal, numbersVal)
			return nil
		},
	}
	infoCmd.Flags().String("pbx-id", "", "PBX identifier to query")
	infoCmd.Flags().String("numbers", "", "Comma-separated list of numbers to filter")

	cmd.AddCommand(infoCmd)

	return cmd
}
