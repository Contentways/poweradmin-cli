// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

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

// NewGetCmd returns a new "groups get" command instance.
func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a group by name or ID",
		Long:  `Get a Poweradmin group by name or numeric ID.`,
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

			var group *poweradmin.Group
			var getErr error
			if idStr != "" {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid id: %w", err)
				}
				group, _, getErr = client.Group.GetByID(cmd.Context(), id)
			} else {
				group, _, getErr = client.Group.GetByName(cmd.Context(), name)
			}
			if getErr != nil {
				return fmt.Errorf("failed to get group: %w", getErr)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(schema.GroupFromSDK(group), "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "ID:          %d\n", group.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "Name:        %s\n", group.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Description: %s\n", group.Description)
			fmt.Fprintf(cmd.OutOrStdout(), "Members:     %d\n", group.MemberCount)
			fmt.Fprintf(cmd.OutOrStdout(), "Zones:       %d\n", group.ZoneCount)

			return nil
		},
	}

	cmd.Flags().String("name", "", "Group name")
	cmd.Flags().String("id", "", "Group ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	return cmd
}
