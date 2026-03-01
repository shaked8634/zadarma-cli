package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/zadarma/zadarma-cli/internal/log"
)

// NewPhoneCmd creates the 'phone' command group for managing phone numbers and virtual numbers.
func NewPhoneCmd(factory ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "phone",
		Short: "Manage phone numbers (DIDs and virtual numbers)",
		Long:  "Commands for managing owned phone numbers and exploring virtual number availability.",
	}

	cmd.AddCommand(
		newPhoneListCmd(factory),
		newPhoneCountriesParentCmd(factory),
		newPhoneCountryParentCmd(factory),
		newPhoneNumberCmd(factory),
	)

	return cmd
}

// newPhoneListCmd returns the 'phone list' command
func newPhoneListCmd(factory ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:   "list [number...]",
		Short: "List phone numbers",
		Long: `List phone numbers (DIDs) owned by your account.

If no numbers are specified, lists all phone numbers.
If one or more numbers are specified, retrieves information for those specific numbers.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			c := factory()

			dids, err := c.GetDirectNumbers(args...)
			if err != nil {
				log.Debugf("Failed to get phone numbers: %v", err)
				return failCmd(cmd, err)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(dids, "", "  ")
				fmt.Println(string(out))
			} else {
				if len(dids) == 0 {
					fmt.Println("No phone numbers found.")
					return nil
				}
				printPhoneTable(dids)
			}
			return nil
		},
	}
}

// newPhoneCountriesParentCmd returns the 'phone countries' parent command which contains subcommands like 'list'
func newPhoneCountriesParentCmd(factory ClientFactory) *cobra.Command {
	parent := &cobra.Command{
		Use:   "countries",
		Short: "return country codes and ISO",
		Long:  "Operations related to available direct number countries (list and other country-level actions)",
	}

	parent.AddCommand(newPhoneCountriesListCmd(factory))
	return parent
}

// newPhoneCountriesListCmd returns the 'phone countries list' command
func newPhoneCountriesListCmd(factory ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List country codes and ISO for direct numbers",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			c := factory()

			countries, err := c.GetDirectCountries()
			if err != nil {
				log.Debugf("Failed to get countries: %v", err)
				return failCmd(cmd, err)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(countries, "", "  ")
				fmt.Println(string(out))
			} else {
				if len(countries) == 0 {
					fmt.Println("No countries found.")
					return nil
				}
				printCountriesTable(countries)
			}
			return nil
		},
	}
}

// newPhoneCountryParentCmd returns the 'phone country' parent command which contains subcommands like 'info'
func newPhoneCountryParentCmd(factory ClientFactory) *cobra.Command {
	parent := &cobra.Command{
		Use:   "country",
		Short: "country-specific operations",
		Long:  "Operations for a specific country, e.g., getting number types in a country",
	}

	parent.AddCommand(newPhoneCountryInfoCmd(factory))
	return parent
}

// newPhoneCountryInfoCmd returns the 'phone country info <code>' command
func newPhoneCountryInfoCmd(factory ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:   "info <code>",
		Short: "return number types per country",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			c := factory()

			countryData, err := c.GetDirectCountry(args[0])
			if err != nil {
				log.Debugf("Failed to get country data: %v", err)
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
				printCountryDestinationsTable(countryData)
			}
			return nil
		},
	}
}

// newPhoneNumberCmd returns the 'phone number' command
func newPhoneNumberCmd(factory ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:   "number <number>",
		Short: "Get information about a virtual number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			c := factory()

			info, err := c.GetDirectNumber(args[0])
			if err != nil {
				log.Debugf("Failed to get number info: %v", err)
				return failCmd(cmd, err)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(info, "", "  ")
				fmt.Println(string(out))
			} else {
				printNumberInfo(info)
			}
			return nil
		},
	}
}

// printPhoneTable renders owned phone numbers as a formatted table
func printPhoneTable(phones []map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() { _ = w.Flush() }()

	// Print header
	_, _ = fmt.Fprintln(w, "NUMBER\tCOUNTRY\tSTATUS\tDESCRIPTION\tEXPIRES")
	_, _ = fmt.Fprintln(w, "------\t-------\t------\t-----------\t-------")

	// Print rows
	for _, phone := range phones {
		number := getStringField(phone, "number", "-")
		country := getStringField(phone, "country", "-")
		status := getStringField(phone, "status", "-")
		description := getStringField(phone, "description", "-")
		stopDate := getStringField(phone, "stop_date", "-")

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", number, country, status, description, stopDate)
	}
}

// printCountriesTable renders countries as a formatted table
func printCountriesTable(countries []map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() { _ = w.Flush() }()

	// Print header
	_, _ = fmt.Fprintln(w, "COUNTRY CODE\tCOUNTRY CODE ISO\tNAME")
	_, _ = fmt.Fprintln(w, "---\t---\t----")

	// Print rows
	for _, country := range countries {
		countryCode := getStringField(country, "countryCode", "-")
		countryCodeISO := getStringField(country, "countryCodeIso", "-")
		name := getStringField(country, "name", "-")

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", countryCode, countryCodeISO, name)
	}
}

// printCountryDestinationsTable renders destinations for a specific country
func printCountryDestinationsTable(destinations []map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() { _ = w.Flush() }()

	_, _ = fmt.Fprintln(w, "ID\tAREA_CODE\tNAME\tMONTHLY_FEE\tRECEIVE_SMS\tIS_TOLL")
	_, _ = fmt.Fprintln(w, "--\t---------\t----\t-----------\t-----------\t-------")

	for _, d := range destinations {
		id := getStringField(d, "id", "-")
		area := getStringField(d, "areaCode", "-")
		name := getStringField(d, "name", "-")
		monthly := getStringField(d, "monthly_fee", "-")
		receiveSMS := getStringField(d, "receive_sms", "-")
		isToll := getStringField(d, "is_toll", "-")

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", id, area, name, monthly, receiveSMS, isToll)
	}
}

// printNumberInfo prints a single direct number's info as a horizontal table
func printNumberInfo(info map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() { _ = w.Flush() }()

	// Deterministic order matching user expectations
	keys := []string{
		"country",
		"currency",
		"channels",
		"autorenew_period",
		"receive_sms",
		"is_on_test",
		"direction_id",
		"number_name",
		"start_date",
		"sip_name",
		"description",
		"stop_date",
		"number",
		"status",
		"monthly_fee",
		"autorenew",
		"group_id",
		"sip",
	}

	// Build header row
	for i, k := range keys {
		if i > 0 {
			_, _ = fmt.Fprint(w, "\t")
		}
		_, _ = fmt.Fprint(w, k)
	}
	_, _ = fmt.Fprint(w, "\n")

	// Build separator row (dashes per header length)
	for i, k := range keys {
		if i > 0 {
			_, _ = fmt.Fprint(w, "\t")
		}
		dashes := strings.Repeat("-", len(k))
		_, _ = fmt.Fprint(w, dashes)
	}
	_, _ = fmt.Fprint(w, "\n")

	// Build value row
	for i, k := range keys {
		if i > 0 {
			_, _ = fmt.Fprint(w, "\t")
		}
		v := getStringField(info, k, "<nil>")
		_, _ = fmt.Fprint(w, v)
	}
	_, _ = fmt.Fprint(w, "\n")
}
