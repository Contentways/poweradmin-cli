// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

package cmd

import (
	"github.com/spf13/cobra"
)

// zonesCmd represents the zones command
var zonesCmd = &cobra.Command{
	Use:   "zones",
	Short: "Manage DNS zones",
	Long:  `Manage DNS zones in Poweradmin — list, create, update and delete zones.`,
}

func init() {
	rootCmd.AddCommand(zonesCmd)
}
