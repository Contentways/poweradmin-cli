// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/internal/output"
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

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			// JSON output — print the raw zone list and return early.
			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(zones, "", "  ")
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
	return cmd
}
