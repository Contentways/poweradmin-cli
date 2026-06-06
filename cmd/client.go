// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/spf13/cobra"
)

func newClient(cmd *cobra.Command) (*poweradmin.Client, error) {
	url, err := cmd.Flags().GetString("url")
	if err != nil || url == "" {
		url, _ = cmd.Root().PersistentFlags().GetString("url")
	}
	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil || apiKey == "" {
		apiKey, _ = cmd.Root().PersistentFlags().GetString("api-key")
	}

	if url == "" {
		return nil, fmt.Errorf("--url is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("--api-key is required")
	}

	return poweradmin.NewClient(
		poweradmin.WithBaseURL(url),
		poweradmin.WithAPIKey(apiKey),
	)
}
