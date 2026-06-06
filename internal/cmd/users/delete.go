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

// NewDeleteCmd returns a new "users delete" command instance.
// A new instance is returned on each call to prevent flag state from leaking
// between successive command executions.
// The user can be identified by username (--name) or numeric ID (--id).
// Output can be a human-readable confirmation (default) or JSON
// containing the deleted user's ID and username.
func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a user",
		Long:  `Delete a Poweradmin user by username or numeric ID.`,
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

			// Resolve user ID — either parse the numeric flag directly,
			// or look up the user by username to obtain the ID.
			var userID int
			var username string
			if idStr != "" {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid id: %w", err)
				}
				userID = id
			} else {
				user, _, err := client.User.GetByName(cmd.Context(), name)
				if err != nil {
					return fmt.Errorf("failed to resolve user: %w", err)
				}
				userID = user.ID
				username = user.Username
			}

			_, err = client.User.Delete(cmd.Context(), userID)
			if err != nil {
				return fmt.Errorf("failed to delete user: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			// JSON output — return the deleted user's ID and username.
			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(map[string]any{
					"id":       userID,
					"username": username,
				}, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			// Default output — human-readable confirmation.
			fmt.Fprintf(cmd.OutOrStdout(), "deleted user %s (id %d)\n", username, userID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Username to identify the user")
	cmd.Flags().String("id", "", "User ID to identify the user")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|full")
	return cmd
}
