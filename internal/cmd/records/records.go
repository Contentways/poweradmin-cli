// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package records

import (
	"github.com/spf13/cobra"
)

func NewRecordsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "records",
		Short: "Manage DNS records",
		Long:  `Manage DNS records in Poweradmin — list, create and delete records.`,
	}

	cmd.AddCommand(ListCmd)
	cmd.AddCommand(CreateCmd)
	cmd.AddCommand(DeleteCmd)

	return cmd
}
