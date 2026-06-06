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

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a DNS record",
	Long:  `Create a new DNS record in a zone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := state.FromContext(cmd.Context())
		client, err := s.Client()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

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
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("created record %s %s %s (id %s)\n", name, recordType, content, id)
		return nil
	},
}

func init() {
	CreateCmd.Flags().String("zone-name", "", "Zone name (e.g. example.com)")
	CreateCmd.Flags().String("zone-id", "", "Zone ID")
	CreateCmd.Flags().String("name", "", "Record name")
	CreateCmd.Flags().String("type", "", "Record type (A, AAAA, CNAME, MX, TXT, NS, ...)")
	CreateCmd.Flags().String("content", "", "Record content")
	CreateCmd.Flags().Int("ttl", 3600, "Time to live in seconds")
	CreateCmd.Flags().Int("priority", 0, "Record priority (for MX records)")
	CreateCmd.Flags().String("output", "table", "Output format: table, json")
}
