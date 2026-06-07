// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package groups

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/contentways/poweradmin-cli/internal/schema"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewListCmd returns a new "groups list" command instance.
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all groups",
		Long:  `List all groups in Poweradmin.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := state.FromContext(cmd.Context())

			client, err := s.Client()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			groups, err := client.Group.All(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to list groups: %w", err)
			}

			outputStr, _ := cmd.Flags().GetString("output")
			outputFmt := output.ParseFormat(outputStr)

			if outputFmt == output.FormatJSON {
				data, err := json.MarshalIndent(schema.GroupListFromSDK(groups), "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			t := output.New(cmd.OutOrStdout())
			t.AddHeader("ID", "NAME", "DESCRIPTION", "MEMBERS", "ZONES")
			noHeader, _ := cmd.Flags().GetBool("no-header")
			t.SetNoHeader(noHeader)
			for _, g := range groups {
				t.AddRow(
					strconv.Itoa(g.ID),
					g.Name,
					g.Description,
					strconv.Itoa(g.MemberCount),
					strconv.Itoa(g.ZoneCount),
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
