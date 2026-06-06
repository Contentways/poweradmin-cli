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

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a DNS record",
	Long:  `Delete a DNS record by ID from a zone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := state.FromContext(cmd.Context())
		client, err := s.Client()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		zoneName, _ := cmd.Flags().GetString("zone-name")
		zoneIDStr, _ := cmd.Flags().GetString("zone-id")
		recordID, _ := cmd.Flags().GetString("id")

		if zoneName == "" && zoneIDStr == "" {
			return fmt.Errorf("either --zone-name or --zone-id is required")
		}
		if recordID == "" {
			return fmt.Errorf("--id is required")
		}

		var zoneID int
		if zoneIDStr != "" {
			zoneID, err = strconv.Atoi(zoneIDStr)
			if err != nil {
				return fmt.Errorf("invalid zone-id: %w", err)
			}
		} else {
			zone, _, err := client.Zone.GetByName(cmd.Context(), zoneName)
			if err != nil {
				return fmt.Errorf("failed to resolve zone: %w", err)
			}
			zoneID = zone.ID
		}

		_, err = client.Record.Delete(cmd.Context(), zoneID, recordID)
		if err != nil {
			return fmt.Errorf("failed to delete record: %w", err)
		}

		outputStr, _ := cmd.Flags().GetString("output")
		outputFmt := output.ParseFormat(outputStr)

		if outputFmt == output.FormatJSON {
			data, err := json.MarshalIndent(map[string]any{
				"id":      recordID,
				"zone_id": zoneID,
			}, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal json: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("deleted record (id %s) from zone (id %d)\n", recordID, zoneID)
		return nil
	},
}

func init() {
	DeleteCmd.Flags().String("zone-name", "", "Zone name (e.g. example.com)")
	DeleteCmd.Flags().String("zone-id", "", "Zone ID")
	DeleteCmd.Flags().String("id", "", "Record ID")
	DeleteCmd.Flags().String("output", "table", "Output format: table, json")
}
