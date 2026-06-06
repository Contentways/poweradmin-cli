// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewDeleteCmd returns a new "records delete" command instance.
// A new instance is returned on each call to prevent flag state from leaking
// between successive command executions.
// The zone can be identified by name (--zone-name) or numeric ID (--zone-id).
// The record ID is the opaque string identifier returned by the Poweradmin API
// in 4.3.0+ API-mode (Base64-encoded JSON).
// Output can be a human-readable confirmation (default) or JSON
// containing the deleted record's ID and zone ID.
func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a DNS record",
		Long:  `Delete a DNS record by ID from a zone.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			zoneName, _ := cmd.Flags().GetString("zone-name")
			zoneIDStr, _ := cmd.Flags().GetString("zone-id")
			recordID, _ := cmd.Flags().GetString("id")

			if zoneName == "" && zoneIDStr == "" {
				return fmt.Errorf("either --zone-name or --zone-id is required")
			}
			if recordID == "" {
				return fmt.Errorf("--id is required")
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			// Resolve zone ID — either parse the numeric flag directly,
			// or look up the zone by name to obtain its ID.
			var zoneID int
			if zoneIDStr != "" {
				zoneID, err = strconv.Atoi(zoneIDStr)
				if err != nil {
					return fmt.Errorf("invalid zone-id: %w", err)
				}
			} else {
				zone, _, err := client.Zone.GetByName(cmd.Context(), zoneName)
				if err != nil {
					return fmt.Errorf("failed to resolve zone: %w", err)
				}
				zoneID = zone.ID
			}

			_, err = client.Record.Delete(cmd.Context(), zoneID, recordID)
			if err != nil {
				return fmt.Errorf("failed to delete record: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			// JSON output — return the deleted record's ID and zone ID.
			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(map[string]any{
					"id":      recordID,
					"zone_id": zoneID,
				}, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			// Default output — human-readable confirmation.
			fmt.Fprintf(cmd.OutOrStdout(), "deleted record (id %s) from zone (id %d)\n", recordID, zoneID)
			return nil
		},
	}

	cmd.Flags().String("zone-name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("zone-id", "", "Zone ID")
	cmd.Flags().String("id", "", "Record ID (opaque string returned by the API)")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|full")
	return cmd
}
