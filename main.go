// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package main

import (
	"fmt"
	"os"

	"github.com/contentways/poweradmin-cli/internal/cli"
	"github.com/contentways/poweradmin-cli/internal/state"
)

func main() {
	url := os.Getenv("POWERADMIN_URL")
	apiKey := os.Getenv("POWERADMIN_API_KEY")

	s, err := state.New(url, apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	root := cli.NewRootCommand(s)
	cli.Execute(root)
}
