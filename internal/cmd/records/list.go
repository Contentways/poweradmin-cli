// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewListCmd returns a new "records list" command instance.
// A new instance is returned on each call to prevent flag state from leaking
// between successive command executions.
// The zone can be identified by name (--zone-name) or numeric ID (--zone-id).
// Output can be formatted as a table (default, content truncated to 50 chars),
// full table (no truncation) or JSON.
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all records in a zone",
		Long:  `List all DNS records in a zone by name or ID.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			name, _ := cmd.Flags().GetString("zone-name")
			idStr, _ := cmd.Flags().GetString("zone-id")

			if name == "" && idStr == "" {
				return fmt.Errorf("either --zone-name or --zone-id is required")
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			// Resolve zone ID — either parse the numeric flag directly,
			// or look up the zone by name to obtain its ID.
			var zoneID int
			if idStr != "" {
				fmt.Sscanf(idStr, "%d", &zoneID)
			} else {
				zone, _, err := client.Zone.GetByName(cmd.Context(), name)
				if err != nil {
					return fmt.Errorf("failed to resolve zone: %w", err)
				}
				zoneID = zone.ID
			}

			records, err := client.Record.All(cmd.Context(), zoneID)
			if err != nil {
				return fmt.Errorf("failed to list records: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			// JSON output — print the full record list and return early.
			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(records, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			// Table output — render an aligned table with NAME, TYPE, CONTENT and TTL.
			// In default table mode, long content values are truncated to 50 characters
			// to keep the output readable. Use --output full to see the complete content.
			t := output.New(cmd.OutOrStdout())
			t.AddHeader("NAME", "TYPE", "CONTENT", "TTL")
			for _, r := range records {
				content := r.Content
				if outputFmt == output.FormatTable {
					content = output.Truncate(content, 50)
				}
				t.AddRow(r.Name, r.Type, content, strconv.Itoa(r.TTL))
			}
			t.Flush()

			return nil
		},
	}

	cmd.Flags().String("zone-name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("zone-id", "", "Zone ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|full")
	return cmd
}
