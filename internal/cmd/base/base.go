// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package base provides shared helpers for CLI commands.
// It reduces boilerplate in command implementations by centralising
// common patterns: JSON output, table creation, delete confirmation,
// zone/user/group resolution by name or ID, and quiet mode.
package base

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/output"
	"github.com/spf13/cobra"
)

// PrintJSON marshals v to indented JSON and writes it to cmd's stdout.
func PrintJSON(cmd *cobra.Command, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}

// NewTable creates a new output.Table writing to cmd's stdout.
// It reads the --no-header flag and configures the table accordingly.
func NewTable(cmd *cobra.Command) *output.Table {
	t := output.New(cmd.OutOrStdout())
	noHeader, _ := cmd.Flags().GetBool("no-header")
	t.SetNoHeader(noHeader)
	return t
}

// IsQuiet returns true if the --quiet flag is set.
func IsQuiet(cmd *cobra.Command) bool {
	q, _ := cmd.Flags().GetBool("quiet")
	return q
}

// Confirm prints a confirmation prompt and reads the user's response.
// Returns true if the user confirmed with "y" or "Y".
// Returns true immediately if the --yes flag is set.
func Confirm(cmd *cobra.Command, msg string) bool {
	yes, _ := cmd.Flags().GetBool("yes")
	if yes {
		return true
	}
	fmt.Fprint(cmd.OutOrStdout(), msg)
	var confirm string
	fmt.Fscan(os.Stdin, &confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
		return false
	}
	return true
}

// ResolveZone returns a Zone identified by --name or --id flag.
// If --id is provided it takes precedence over --name.
func ResolveZone(cmd *cobra.Command, client *poweradmin.Client) (*poweradmin.Zone, error) {
	name, _ := cmd.Flags().GetString("name")
	idStr, _ := cmd.Flags().GetString("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %w", err)
		}
		zone, _, err := client.Zone.GetByID(cmd.Context(), id)
		if err != nil {
			return nil, fmt.Errorf("failed to get zone: %w", err)
		}
		return zone, nil
	}

	zone, _, err := client.Zone.GetByName(cmd.Context(), name)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve zone: %w", err)
	}
	return zone, nil
}

// ResolveUser returns a User identified by --name or --id flag.
func ResolveUser(cmd *cobra.Command, client *poweradmin.Client) (*poweradmin.User, error) {
	name, _ := cmd.Flags().GetString("name")
	idStr, _ := cmd.Flags().GetString("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %w", err)
		}
		user, _, err := client.User.GetByID(cmd.Context(), id)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
		return user, nil
	}

	user, _, err := client.User.GetByName(cmd.Context(), name)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve user: %w", err)
	}
	return user, nil
}

// ResolveGroup returns a Group identified by --name or --id flag.
func ResolveGroup(cmd *cobra.Command, client *poweradmin.Client) (*poweradmin.Group, error) {
	name, _ := cmd.Flags().GetString("name")
	idStr, _ := cmd.Flags().GetString("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %w", err)
		}
		group, _, err := client.Group.GetByID(cmd.Context(), id)
		if err != nil {
			return nil, fmt.Errorf("failed to get group: %w", err)
		}
		return group, nil
	}

	group, _, err := client.Group.GetByName(cmd.Context(), name)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve group: %w", err)
	}
	return group, nil
}

// ResolveZoneID returns the numeric zone ID from --zone-name or --zone-id flag.
func ResolveZoneID(cmd *cobra.Command, client *poweradmin.Client) (int, error) {
	zoneName, _ := cmd.Flags().GetString("zone-name")
	zoneIDStr, _ := cmd.Flags().GetString("zone-id")

	if zoneIDStr != "" {
		id, err := strconv.Atoi(zoneIDStr)
		if err != nil {
			return 0, fmt.Errorf("invalid zone-id: %w", err)
		}
		return id, nil
	}

	zone, _, err := client.Zone.GetByName(cmd.Context(), zoneName)
	if err != nil {
		return 0, fmt.Errorf("failed to resolve zone: %w", err)
	}
	return zone.ID, nil
}
