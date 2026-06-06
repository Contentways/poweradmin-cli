// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users

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

// NewGetCmd returns a new "users get" command instance.
// A new instance is returned on each call to prevent flag state from leaking
// between successive command executions.
// The user can be identified by username (--name) or numeric ID (--id).
// Output can be a human-readable key-value summary (default) or JSON.
func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a user by name or ID",
		Long:  `Get a Poweradmin user by username or numeric ID.`,
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

			// Resolve the user — either by numeric ID or by username.
			var user *poweradmin.User
			var getErr error
			if idStr != "" {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid id: %w", err)
				}
				user, _, getErr = client.User.GetByID(cmd.Context(), id)
			} else {
				user, _, getErr = client.User.GetByName(cmd.Context(), name)
			}
			if getErr != nil {
				return fmt.Errorf("failed to get user: %w", getErr)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			// JSON output — return the full user object.
			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(schema.UserFromSDK(user), "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			// Default output — human-readable key-value summary.
			active := "no"
			if user.Active {
				active = "yes"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "ID:       %d\n", user.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "Username: %s\n", user.Username)
			fmt.Fprintf(cmd.OutOrStdout(), "Email:    %s\n", user.Email)
			fmt.Fprintf(cmd.OutOrStdout(), "Active:   %s\n", active)
			if user.Fullname != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Fullname: %s\n", user.Fullname)
			}

			return nil
		},
	}

	cmd.Flags().String("name", "", "Username")
	cmd.Flags().String("id", "", "User ID")
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|full")
	return cmd
}
