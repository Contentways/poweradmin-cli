// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package cli

import (
	"os"

	"github.com/contentways/poweradmin-cli/internal/cmd/records"
	"github.com/contentways/poweradmin-cli/internal/cmd/zones"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

func NewRootCommand(s *state.State) *cobra.Command {
	root := &cobra.Command{
		Use:   "poweradmin",
		Short: "CLI for managing Poweradmin DNS",
		Long:  `poweradmin is a command-line tool for managing DNS zones and records via the Poweradmin REST API.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := s.WithContext(cmd.Context())
			cmd.SetContext(ctx)
			return nil
		},
	}

	root.PersistentFlags().String("url", "", "Poweradmin URL (e.g. https://dns.example.com)")
	root.PersistentFlags().String("api-key", "", "Poweradmin API key")

	root.AddCommand(zones.NewZonesCommand())
	root.AddCommand(records.NewRecordsCommand())

	return root
}

func Execute(root *cobra.Command) {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
