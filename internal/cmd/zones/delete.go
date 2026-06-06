// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

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
		}

		_, err = client.Zone.Delete(cmd.Context(), zoneID)
		if err != nil {
			return fmt.Errorf("failed to delete zone: %w", err)
		}

		fmt.Printf("deleted zone (id %d)\n", zoneID)
		return nil
	},
}

func init() {
	DeleteCmd.Flags().String("name", "", "Zone name (e.g. example.com)")
	DeleteCmd.Flags().String("id", "", "Zone ID")
}
