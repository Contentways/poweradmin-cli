// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package groups provides CLI commands for managing groups in Poweradmin.
// Available subcommands: list, get, create, update, delete,
// members, member-add, member-remove, zones, zone-add, zone-remove.
package groups

import (
	"github.com/spf13/cobra"
)

// NewGroupsCommand builds and returns the "groups" subcommand group.
// All group-related commands are registered here and made available
// under the "poweradmin groups" namespace.
func NewGroupsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups",
		Short: "Manage Poweradmin groups",
		Long:  `Manage Poweradmin groups — list, get, create, update, delete and manage members and zones.`,
	}

	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewGetCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewDeleteCmd())
	cmd.AddCommand(NewMembersCmd())
	cmd.AddCommand(NewMemberAddCmd())
	cmd.AddCommand(NewMemberRemoveCmd())
	cmd.AddCommand(NewZonesCmd())
	cmd.AddCommand(NewZoneAddCmd())
	cmd.AddCommand(NewZoneRemoveCmd())

	return cmd
}
