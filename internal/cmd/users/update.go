// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users

import (
	"encoding/json"
	"fmt"
	"strconv"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewUpdateCmd returns a new "users update" command instance.
// A new instance is returned on each call to prevent flag state from leaking
// between successive command executions.
// The user can be identified by username (--name) or numeric ID (--id).
// Only flags that are explicitly set are sent to the API — unset flags are omitted.
// Output can be a human-readable confirmation (default) or JSON.
func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a user",
		Long:  `Update an existing Poweradmin user by username or numeric ID.`,
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
			}

			// Build update opts — only include fields that were explicitly set.
			opts := poweradmin.UserUpdateOpts{}
			if cmd.Flags().Changed("email") {
				opts.Email, _ = cmd.Flags().GetString("email")
			}
			if cmd.Flags().Changed("fullname") {
				opts.Fullname, _ = cmd.Flags().GetString("fullname")
			}
			if cmd.Flags().Changed("password") {
				opts.Password, _ = cmd.Flags().GetString("password")
			}
			if cmd.Flags().Changed("active") {
				active, _ := cmd.Flags().GetBool("active")
				opts.Active = &active
			}

			user, _, err := client.User.Update(cmd.Context(), userID, opts)
			if err != nil {
				return fmt.Errorf("failed to update user: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			// JSON output — return the updated user object.
			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(user, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			// Default output — human-readable confirmation.
			fmt.Fprintf(cmd.OutOrStdout(), "updated user %s (id %d)\n", user.Username, user.ID)
			return nil
		},
	}

	cmd.Flags().String("name", "", "Username to identify the user")
	cmd.Flags().String("id", "", "User ID to identify the user")
	cmd.Flags().String("email", "", "New email address")
	cmd.Flags().String("fullname", "", "New full name")
	cmd.Flags().String("password", "", "New password")
	cmd.Flags().Bool("active", true, "Whether the user is active")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|full")
	return cmd
}
