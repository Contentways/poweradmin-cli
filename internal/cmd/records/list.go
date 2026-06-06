// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"fmt"

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

		for _, r := range records {
			fmt.Printf("%s\t%s\t%s\tTTL=%d\n", r.Name, r.Type, r.Content, r.TTL)
		}

		return nil
	},
}

func init() {
	ListCmd.Flags().String("zone-name", "", "Zone name (e.g. example.com)")
	ListCmd.Flags().String("zone-id", "", "Zone ID")
}
