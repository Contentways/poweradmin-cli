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

// NewCreateCmd returns a new "zones create" command instance.
// A new instance is returned on each call to prevent flag state from leaking
// between successive command executions.
// The zone name is passed as a positional argument.
// Output can be a human-readable confirmation (default) or JSON
// containing the new zone's ID, name and type.
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
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

			// JSON output — return the new zone's ID, name and type.
			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(map[string]any{
					"id":   id,
					"name": args[0],
					"type": zoneType,
				}, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			// Default output — human-readable confirmation.
			fmt.Fprintf(cmd.OutOrStdout(), "created zone %s (id %d)\n", args[0], id)
			return nil
		},
	}

	cmd.Flags().String("type", "NATIVE", "Zone type. One of: NATIVE|MASTER|SLAVE")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|full")
	return cmd
}
