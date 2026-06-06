// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"encoding/json"
	"fmt"
	"strconv"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/schema"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewGetCmd returns a new "zones get" command instance.
// A new instance is returned on each call to prevent flag state from leaking
// between successive command executions.
// The zone can be identified by name (--name) or numeric ID (--id).
// Nameservers are resolved by fetching all NS records for the zone.
func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
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

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			// Resolve the zone — either by numeric ID or by name.
			var zone *poweradmin.Zone
			var getErr error
			if idStr != "" {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid id: %w", err)
				}
				zone, _, getErr = client.Zone.GetByID(cmd.Context(), id)
			} else {
				zone, _, getErr = client.Zone.GetByName(cmd.Context(), name)
			}
			if getErr != nil {
				return fmt.Errorf("failed to get zone: %w", getErr)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			// JSON output — fetch records, extract NS entries, and return a combined object.
			if outputFmt == output.FormatJSON {
				records, err := client.Record.All(cmd.Context(), zone.ID)
				if err != nil {
					return fmt.Errorf("failed to get records: %w", err)
				}

				var nameservers []string
				for _, r := range records {
					if r.Type == "NS" {
						nameservers = append(nameservers, r.Content)
					}
				}

				data, err := json.MarshalIndent(schema.ZoneWithNameservers{
					Zone:        schema.ZoneFromSDK(zone),
					Nameservers: nameservers,
				}, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			// Default output — key-value summary with nameservers.
			records, err := client.Record.All(cmd.Context(), zone.ID)
			if err != nil {
				return fmt.Errorf("failed to get records: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "ID:   %d\n", zone.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "Name: %s\n", zone.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Type: %s\n", zone.Type)
			if zone.Masters != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Masters: %s\n", zone.Masters)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Nameservers:")
			for _, r := range records {
				if r.Type == "NS" {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", r.Content)
				}
			}

			return nil
		},
	}

	cmd.Flags().String("name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("id", "", "Zone ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|full")
	return cmd
}
