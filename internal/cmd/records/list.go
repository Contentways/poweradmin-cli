// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all records in a zone",
	Long:  `List all DNS records in a zone by name or ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := state.FromContext(cmd.Context())
		client, err := s.Client()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		name, _ := cmd.Flags().GetString("zone-name")
		idStr, _ := cmd.Flags().GetString("zone-id")

		if name == "" && idStr == "" {
			return fmt.Errorf("either --zone-name or --zone-id is required")
		}

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

		if outputFmt == output.FormatJSON {
			data, err := json.MarshalIndent(records, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal json: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		t := output.New(os.Stdout)
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

func init() {
	ListCmd.Flags().String("zone-name", "", "Zone name (e.g. example.com)")
	ListCmd.Flags().String("zone-id", "", "Zone ID")
	ListCmd.Flags().String("output", "table", "Output format: table, full, json")
}
