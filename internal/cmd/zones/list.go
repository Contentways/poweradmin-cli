// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"fmt"

	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all DNS zones",
	Long:  `List all DNS zones in Poweradmin.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := state.FromContext(cmd.Context())

		zones, err := s.Client().Zone.All(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to list zones: %w", err)
		}

		for _, z := range zones {
			fmt.Printf("%d\t%s\t%s\n", z.ID, z.Name, z.Type)
		}

		return nil
	},
}
