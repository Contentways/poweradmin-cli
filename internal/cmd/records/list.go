// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/schema"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewListCmd returns a new "records list" command instance.
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all records in a zone",
		Long:  `List all DNS records in a zone by name or ID.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			zoneName, _ := cmd.Flags().GetString("zone-name")
			zoneIDStr, _ := cmd.Flags().GetString("zone-id")

			if zoneName == "" && zoneIDStr == "" {
				return fmt.Errorf("either --zone-name or --zone-id is required")
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			zoneID, err := base.ResolveZoneID(cmd, client)
			if err != nil {
				return err
			}

			records, err := client.Record.All(cmd.Context(), zoneID)
			if err != nil {
				return fmt.Errorf("failed to list records: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, schema.RecordListFromSDK(records))
			}

			t := base.NewTable(cmd)
			t.AddHeader("NAME", "TYPE", "CONTENT", "TTL")
			for _, r := range records {
				content := r.Content
				if outputFmt == output.FormatTable {
					content = output.Truncate(content, 50)
				}
				t.AddColoredRow(
					output.PlainCell(r.Name),
					output.Cell(r.Type, output.CyanCode()),
					output.PlainCell(content),
					output.PlainCell(strconv.Itoa(r.TTL)),
				)
			}
			t.Flush()
			return nil
		},
	}

	cmd.Flags().String("zone-name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("zone-id", "", "Zone ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|full|json")
	cmd.Flags().Bool("no-header", false, "Suppress table header row")
	return cmd
}
