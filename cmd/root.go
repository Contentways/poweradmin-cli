// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "poweradmin",
	Short: "CLI for managing Poweradmin DNS",
	Long:  `poweradmin is a command-line tool for managing DNS zones and records via the Poweradmin REST API.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("url", "", "Poweradmin URL (e.g. https://dns.example.com)")
	rootCmd.PersistentFlags().String("api-key", "", "Poweradmin API key")
}
