// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package records provides CLI commands for managing DNS records in Poweradmin.
// Available subcommands: list, create, delete.
package records

import (
	"github.com/spf13/cobra"
)

// NewRecordsCommand builds and returns the "records" subcommand group.
// All record-related commands are registered here and made available
// under the "poweradmin records" namespace.
func NewRecordsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "records",
		Short: "Manage DNS records",
		Long:  `Manage DNS records in Poweradmin — list, create and delete records.`,
	}

	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewDeleteCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewGetCmd())

	return cmd
}
