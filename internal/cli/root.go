// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package cli assembles the root Cobra command and registers all subcommands.
package cli

import (
	"os"

	"github.com/contentways/poweradmin-cli/internal/cmd/groups"
	"github.com/contentways/poweradmin-cli/internal/cmd/records"
	"github.com/contentways/poweradmin-cli/internal/cmd/users"
	"github.com/contentways/poweradmin-cli/internal/cmd/zones"
	"github.com/contentways/poweradmin-cli/internal/state"
	"github.com/spf13/cobra"
)

// NewRootCommand builds and returns the root Cobra command.
// It registers all persistent flags and subcommands, and sets up
// PersistentPreRunE to inject the State into the command context before
// any subcommand runs. This makes the State available to all commands
// via state.FromContext(cmd.Context()).
func NewRootCommand(s *state.State) *cobra.Command {
	root := &cobra.Command{
		Use:   "poweradmin",
		Short: "CLI for managing Poweradmin DNS",
		Long:  `poweradmin is a command-line tool for managing DNS zones and records via the Poweradmin REST API.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// CLI flags take the highest precedence — override any values
			// that were loaded from the config file or environment variables.
			if url, _ := cmd.Root().PersistentFlags().GetString("url"); url != "" {
				s.URL = url
			}
			if apiKey, _ := cmd.Root().PersistentFlags().GetString("api-key"); apiKey != "" {
				s.APIKey = apiKey
			}

			// Embed the State into the context so all subcommands can
			// retrieve it via state.FromContext(cmd.Context()).
			ctx := s.WithContext(cmd.Context())
			cmd.SetContext(ctx)
			return nil
		},
	}

	// Persistent flags are available to the root command and all subcommands.
	// They override values from the config file and environment variables.
	root.PersistentFlags().StringP("url", "u", "", "Poweradmin URL (e.g. https://dns.example.com)")
	root.PersistentFlags().StringP("api-key", "k", "", "Poweradmin API key (overrides config and env)")

	// Register all resource subcommands.
	root.AddCommand(zones.NewZonesCommand())
	root.AddCommand(records.NewRecordsCommand())
	root.AddCommand(users.NewUsersCommand())
	root.AddCommand(groups.NewGroupsCommand())

	return root
}

// Execute runs the root command and exits with a non-zero status code on error.
// This is the single entry point called from main.
func Execute(root *cobra.Command) {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
