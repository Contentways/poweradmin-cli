// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var zonesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all DNS zones",
	Long:  `List all DNS zones in Poweradmin.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient(cmd)
		if err != nil {
			return err
		}

		zones, err := client.Zone.All(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to list zones: %w", err)
		}

		for _, z := range zones {
			fmt.Printf("%d\t%s\t%s\n", z.ID, z.Name, z.Type)
		}

		return nil
	},
}

func init() {
	zonesCmd.AddCommand(zonesListCmd)
}
