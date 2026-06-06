// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones

import (
	"github.com/spf13/cobra"
)

func NewZonesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zones",
		Short: "Manage DNS zones",
		Long:  `Manage DNS zones in Poweradmin — list, get, create and delete zones.`,
	}

	cmd.AddCommand(ListCmd)
	cmd.AddCommand(GetCmd)
	cmd.AddCommand(CreateCmd)
	cmd.AddCommand(DeleteCmd)

	return cmd
}
