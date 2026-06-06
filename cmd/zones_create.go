// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package cmd

import (
	"fmt"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/spf13/cobra"
)

var zonesCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a DNS zone",
	Long:  `Create a new DNS zone in Poweradmin.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient(cmd)
		if err != nil {
			return err
		}

		zoneType, _ := cmd.Flags().GetString("type")

		id, _, err := client.Zone.Create(cmd.Context(), poweradmin.ZoneCreateOpts{
			Name: args[0],
			Type: poweradmin.ZoneType(zoneType),
		})
		if err != nil {
			return fmt.Errorf("failed to create zone: %w", err)
		}

		fmt.Printf("created zone %s (id %d)\n", args[0], id)
		return nil
	},
}

func init() {
	zonesCreateCmd.Flags().String("type", "NATIVE", "Zone type: NATIVE, MASTER or SLAVE")
	zonesCmd.AddCommand(zonesCreateCmd)
}
