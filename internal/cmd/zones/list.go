// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/schema"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewListCmd returns a new "zones list" command instance.
// A new instance is returned on each call to prevent flag state from leaking
// between successive command executions — this is especially important in tests
// where the same command may be run multiple times.
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all DNS zones",
		Long:  `List all DNS zones in Poweradmin.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			zones, err := client.Zone.All(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to list zones: %w", err)
			}

			// Apply optional filters client-side.
			typeFilter, _ := cmd.Flags().GetString("type")
			nameFilter, _ := cmd.Flags().GetString("name-filter")

			if typeFilter != "" || nameFilter != "" {
				filtered := zones[:0]
				for _, z := range zones {
					if typeFilter != "" && !strings.EqualFold(string(z.Type), typeFilter) {
						continue
					}
					if nameFilter != "" && !strings.Contains(z.Name, nameFilter) {
						continue
					}
					filtered = append(filtered, z)
				}
				zones = filtered
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			// JSON output — print zones wrapped in a root object.
			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(schema.ZoneListFromSDK(zones), "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			// Table output — render an aligned table with ID, name and type columns.
			t := output.New(cmd.OutOrStdout())
			t.AddHeader("ID", "NAME", "TYPE")
			for _, z := range zones {
				t.AddRow(strconv.Itoa(z.ID), z.Name, string(z.Type))
			}
			t.Flush()

			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|full")
	cmd.Flags().String("type", "", "Filter by zone type. One of: NATIVE|MASTER|SLAVE")
	cmd.Flags().String("name-filter", "", "Filter by zone name (substring match)")
	return cmd
}
