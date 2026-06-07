// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewDeleteCmd returns a new "groups delete" command instance.
func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a group",
		Long:  `Delete a Poweradmin group by name or numeric ID.`,
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

			var groupID int
			var groupName string
			if idStr != "" {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid id: %w", err)
				}
				groupID = id
			} else {
				group, _, err := client.Group.GetByName(cmd.Context(), name)
				if err != nil {
					return fmt.Errorf("failed to resolve group: %w", err)
				}
				groupID = group.ID
				groupName = group.Name
			}

			// Confirm deletion unless --yes flag is set.
			yes, _ := cmd.Flags().GetBool("yes")
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Delete group (id %d)? [y/N] ", groupID)
				var confirm string
				fmt.Fscan(os.Stdin, &confirm)
				if confirm != "y" && confirm != "Y" {
					fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
					return nil
				}
			}

			_, err = client.Group.Delete(cmd.Context(), groupID)
			if err != nil {
				return fmt.Errorf("failed to delete group: %w", err)
			}

			quiet, _ := cmd.Flags().GetBool("quiet")
			if quiet {
				return nil
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(map[string]any{
					"id":   groupID,
					"name": groupName,
				}, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted group %s (id %d)\n", groupName, groupID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Group name to identify the group")
	cmd.Flags().String("id", "", "Group ID to identify the group")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolP("quiet", "q", false, "Suppress output after deletion")

	return cmd
}
