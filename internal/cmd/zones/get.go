// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"fmt"
	"strconv"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a DNS zone by name or ID",
	Long:  `Get a DNS zone by name or ID from Poweradmin.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := state.FromContext(cmd.Context())

		name, _ := cmd.Flags().GetString("name")
		idStr, _ := cmd.Flags().GetString("id")

		if name == "" && idStr == "" {
			return fmt.Errorf("either --name or --id is required")
		}

		var zone *poweradmin.Zone
		var err error

		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}
			zone, _, err = s.Client().Zone.GetByID(cmd.Context(), id)
		} else {
			zone, _, err = s.Client().Zone.GetByName(cmd.Context(), name)
		}
		if err != nil {
			return fmt.Errorf("failed to get zone: %w", err)
		}

		records, err := s.Client().Record.All(cmd.Context(), zone.ID)
		if err != nil {
			return fmt.Errorf("failed to get records: %w", err)
		}

		fmt.Printf("ID:   %d\n", zone.ID)
		fmt.Printf("Name: %s\n", zone.Name)
		fmt.Printf("Type: %s\n", zone.Type)
		if zone.Masters != "" {
			fmt.Printf("Masters: %s\n", zone.Masters)
		}

		fmt.Println("Nameservers:")
		for _, r := range records {
			if r.Type == "NS" {
				fmt.Printf("  %s\n", r.Content)
			}
		}

		return nil
	},
}

func init() {
	GetCmd.Flags().String("name", "", "Zone name (e.g. example.com)")
	GetCmd.Flags().String("id", "", "Zone ID")
}
