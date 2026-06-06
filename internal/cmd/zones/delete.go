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

// DeleteCmd deletes a DNS zone from Poweradmin by name or numeric ID.
// If a name is provided, it is first resolved to an ID via the API.
// Output can be a human-readable confirmation (default) or JSON
// containing the deleted zone's ID and name.
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a DNS zone",
	Long:  `Delete a DNS zone from Poweradmin by name or ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := state.FromContext(cmd.Context())

		name, _ := cmd.Flags().GetString("name")
		idStr, _ := cmd.Flags().GetString("id")

		client, err := s.Client()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if name == "" && idStr == "" {
			return fmt.Errorf("either --name or --id is required")
		}

		// Resolve zone ID — either parse the numeric flag directly,
		// or look up the zone by name to obtain its ID.
		var zoneID int
		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}
			zoneID = id
		} else {
			zone, _, err := client.Zone.GetByName(cmd.Context(), name)
			if err != nil {
				return fmt.Errorf("failed to resolve zone: %w", err)
			}
			zoneID = zone.ID
			name = zone.Name
		}

		_, err = client.Zone.Delete(cmd.Context(), zoneID)
		if err != nil {
			return fmt.Errorf("failed to delete zone: %w", err)
		}

		outputStr, _ := cmd.Flags().GetString("output")
		outputFmt := output.ParseFormat(outputStr)

		// JSON output — return the deleted zone's ID and name.
		if outputFmt == output.FormatJSON {
			data, err := json.MarshalIndent(map[string]any{
				"id":   zoneID,
				"name": name,
			}, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal json: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Default output — human-readable confirmation.
		fmt.Printf("deleted zone %s (id %d)\n", name, zoneID)
		return nil
	},
}

func init() {
	DeleteCmd.Flags().String("name", "", "Zone name (e.g. example.com)")
	DeleteCmd.Flags().String("id", "", "Zone ID")
	DeleteCmd.Flags().String("output", "table", "Output format. One of: table|json")
}
