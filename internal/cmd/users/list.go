// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package users

import (
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/schema"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewListCmd returns a new "users list" command instance.
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

			if outputFmt == output.FormatJSON {
				return base.PrintJSON(cmd, schema.UserListFromSDK(users))
			}

			t := base.NewTable(cmd)
			t.AddHeader("ID", "USERNAME", "EMAIL", "ACTIVE")
			for _, u := range users {
				active := "no"
				activeColor := output.RedCode()
				if u.Active {
					active = "yes"
					activeColor = output.GreenCode()
				}
				t.AddColoredRow(
					output.PlainCell(strconv.Itoa(u.ID)),
					output.PlainCell(u.Username),
					output.PlainCell(u.Email),
					output.Cell(active, activeColor),
				)
			}
			t.Flush()
			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "table", "Output format. One of: table|json")
	cmd.Flags().Bool("no-header", false, "Suppress table header row")
	return cmd
}
