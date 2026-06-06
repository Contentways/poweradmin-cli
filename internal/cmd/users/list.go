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

// NewListCmd returns a new "users list" command instance.
// A new instance is returned on each call to prevent flag state from leaking
// between successive command executions.
// Output can be formatted as a table (default) or JSON.
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all users",
		Long:  `List all users in Poweradmin.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			users, err := client.User.All(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to list users: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			// JSON output — print the raw user list and return early.
			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(users, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			// Table output — render an aligned table with ID, username, email and active columns.
			t := output.New(cmd.OutOrStdout())
			t.AddHeader("ID", "USERNAME", "EMAIL", "ACTIVE")
			for _, u := range users {
				active := "no"
				if u.Active {
					active = "yes"
				}
				t.AddRow(strconv.Itoa(u.ID), u.Username, u.Email, active)
			}
			t.Flush()

			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json|full")
	return cmd
}
