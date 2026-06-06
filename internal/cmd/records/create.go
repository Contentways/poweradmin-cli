// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"encoding/json"
	"fmt"
	"strconv"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewCreateCmd returns a new "records create" command instance.
// A new instance is returned on each call to prevent flag state from leaking
// between successive command executions.
// The zone can be identified by name (--zone-name) or numeric ID (--zone-id).
// Output can be a human-readable confirmation (default) or JSON
// containing the new record's ID and all input fields.
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a DNS record",
		Long:  `Create a new DNS record in a zone.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			zoneName, _ := cmd.Flags().GetString("zone-name")
			zoneIDStr, _ := cmd.Flags().GetString("zone-id")
			name, _ := cmd.Flags().GetString("name")
			recordType, _ := cmd.Flags().GetString("type")
			content, _ := cmd.Flags().GetString("content")
			ttl, _ := cmd.Flags().GetInt("ttl")
			priority, _ := cmd.Flags().GetInt("priority")

			if zoneName == "" && zoneIDStr == "" {
				return fmt.Errorf("either --zone-name or --zone-id is required")
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

			id, _, err := client.Record.Create(cmd.Context(), zoneID, poweradmin.RecordCreateOpts{
				Name:     name,
				Type:     recordType,
				Content:  content,
				TTL:      ttl,
				Priority: priority,
			})
			if err != nil {
				return fmt.Errorf("failed to create record: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			// JSON output — return the new record's ID and all input fields.
			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(map[string]any{
					"id":      id,
					"name":    name,
					"type":    recordType,
					"content": content,
					"ttl":     ttl,
					"zone_id": zoneID,
				}, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			// Default output — human-readable confirmation.
			fmt.Fprintf(cmd.OutOrStdout(), "created record %s %s %s (id %s)\n", name, recordType, content, id)
			return nil
		},
	}

	cmd.Flags().String("zone-name", "", "Zone name (e.g. example.com)")
	cmd.Flags().String("zone-id", "", "Zone ID")
	cmd.Flags().String("name", "", "Record name (e.g. www.example.com)")
	cmd.Flags().String("type", "", "Record type. One of: A|AAAA|CNAME|MX|TXT|NS|SRV|...")
	cmd.Flags().String("content", "", "Record content (e.g. 1.2.3.4 for A records)")
	cmd.Flags().Int("ttl", 3600, "Time to live in seconds (default: 3600)")
	cmd.Flags().Int("priority", 0, "Record priority, used for MX records (default: 0)")
	cmd.Flags().String("output", "table", "Output format. One of: table|json")
	return cmd
}
