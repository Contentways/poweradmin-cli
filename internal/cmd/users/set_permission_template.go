// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewSetPermissionTemplateCmd returns a new "users set-permission-template" command instance.
// A new instance is returned on each call to prevent flag state from leaking
// between successive command executions.
// The user can be identified by username (--name) or numeric ID (--id).
// The permission template is identified by its numeric ID (--template-id).
// Output can be a human-readable confirmation (default) or JSON.
func NewSetPermissionTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-permission-template",
		Short: "Assign a permission template to a user",
		Long:  `Assign a permission template to a Poweradmin user by username or numeric ID.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			name, _ := cmd.Flags().GetString("name")
			idStr, _ := cmd.Flags().GetString("id")
			templateIDStr, _ := cmd.Flags().GetString("template-id")

			if name == "" && idStr == "" {
				return fmt.Errorf("either --name or --id is required")
			}
			if templateIDStr == "" {
				return fmt.Errorf("--template-id is required")
			}

			templateID, err := strconv.Atoi(templateIDStr)
			if err != nil {
				return fmt.Errorf("invalid template-id: %w", err)
			}

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			// Resolve user ID — either parse the numeric flag directly,
			// or look up the user by username to obtain the ID.
			var userID int
			var username string
			if idStr != "" {
				userID, err = strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid id: %w", err)
				}
			} else {
				user, _, err := client.User.GetByName(cmd.Context(), name)
				if err != nil {
					return fmt.Errorf("failed to resolve user: %w", err)
				}
				userID = user.ID
				username = user.Username
			}

			_, err = client.User.SetPermissionTemplate(cmd.Context(), userID, templateID)
			if err != nil {
				return fmt.Errorf("failed to set permission template: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			// JSON output — return the user ID and assigned template ID.
			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(map[string]any{
					"user_id":     userID,
					"username":    username,
					"template_id": templateID,
				}, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			// Default output — human-readable confirmation.
			fmt.Fprintf(cmd.OutOrStdout(), "assigned permission template %d to user %s (id %d)\n", templateID, username, userID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Username to identify the user")
	cmd.Flags().String("id", "", "User ID to identify the user")
	cmd.Flags().String("template-id", "", "Permission template ID to assign (required)")
	cmd.Flags().String("output", "table", "Output format. One of: table|json")
	return cmd
}
