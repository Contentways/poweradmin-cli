// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"encoding/json"
	"fmt"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a DNS zone",
	Long:  `Create a new DNS zone in Poweradmin.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s := state.FromContext(cmd.Context())
		zoneType, _ := cmd.Flags().GetString("type")
		client, err := s.Client()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		id, _, err := client.Zone.Create(cmd.Context(), poweradmin.ZoneCreateOpts{
			Name: args[0],
			Type: poweradmin.ZoneType(zoneType),
		})
		if err != nil {
			return fmt.Errorf("failed to create zone: %w", err)
		}

		outputStr, _ := cmd.Flags().GetString("output")
		outputFmt := output.ParseFormat(outputStr)

		if outputFmt == output.FormatJSON {
			data, err := json.MarshalIndent(map[string]any{
				"id":   id,
				"name": args[0],
				"type": zoneType,
			}, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal json: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("created zone %s (id %d)\n", args[0], id)
		return nil
	},
}

func init() {
	CreateCmd.Flags().String("type", "NATIVE", "Zone type: NATIVE, MASTER or SLAVE")
	CreateCmd.Flags().String("output", "table", "Output format: table, json")
}
