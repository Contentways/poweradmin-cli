// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

//go:build tools

// Package tools tracks tool dependencies that are not imported by the main
// module but are required at build time (e.g. doc generation).
package tools

import _ "github.com/spf13/cobra/doc"
