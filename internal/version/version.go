// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package version holds build-time version information injected via ldflags.
package version

// Version and Commit are set at build time via GoReleaser ldflags.
// Defaults to "dev" and "none" for local builds.
var (
	Version = "dev"
	Commit  = "none"
)
