// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package zones_test

import (
	"context"
	"strings"
	"testing"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/cmd/zones"
	"github.com/contentways/poweradmin-cli/internal/testutil"
)

func TestZonesDelete(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 42, Name: name, Type: "NATIVE"}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.DeleteCmd, []string{"--name", "example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, "example.com") {
		t.Errorf("expected output to contain example.com, got:\n%s", out)
	}
	if !strings.Contains(out, "42") {
		t.Errorf("expected output to contain id 42, got:\n%s", out)
	}
}

func TestZonesDeleteJSON(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 42, Name: name, Type: "NATIVE"}, nil, nil
		},
		DeleteFn: func(ctx context.Context, id int) (*poweradmin.Response, error) {
			return nil, nil
		},
	}

	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	err := fx.Run(zones.DeleteCmd, []string{"--name", "example.com", "--output", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := fx.Stdout.String()
	if !strings.Contains(out, `"name": "example.com"`) {
		t.Errorf("expected JSON to contain example.com, got:\n%s", out)
	}
	if !strings.Contains(out, `"id": 42`) {
		t.Errorf("expected JSON to contain id 42, got:\n%s", out)
	}
}
