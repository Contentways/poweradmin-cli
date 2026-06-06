// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var zonesGetCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Get a DNS zone by name",
	Long:  `Get a DNS zone by name from Poweradmin.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient(cmd)
		if err != nil {
			return err
		}

		zone, _, err := client.Zone.GetByName(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to get zone: %w", err)
		}

		records, err := client.Record.All(cmd.Context(), zone.ID)
		if err != nil {
			return fmt.Errorf("failed to get records: %w", err)
		}

		fmt.Printf("ID:      %d\n", zone.ID)
		fmt.Printf("Name:    %s\n", zone.Name)
		fmt.Printf("Type:    %s\n", zone.Type)
		if zone.Masters != "" {
			fmt.Printf("Masters: %s\n", zone.Masters)
		}

		fmt.Printf("Nameservers:\n")
		for _, r := range records {
			if r.Type == "NS" {
				fmt.Printf("  %s\n", r.Content)
			}
		}

		return nil
	},
}

func init() {
	zonesCmd.AddCommand(zonesGetCmd)
}
