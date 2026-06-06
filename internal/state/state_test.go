// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package state_test

import (
	"context"
	"testing"

	"github.com/contentways/poweradmin-cli/internal/state"
)

func TestNew(t *testing.T) {
	s := state.New("https://dns.example.com", "pwa_test")
	if s.URL != "https://dns.example.com" {
		t.Errorf("URL = %q, want https://dns.example.com", s.URL)
	}
	if s.APIKey != "pwa_test" {
		t.Errorf("APIKey = %q, want pwa_test", s.APIKey)
	}
}

func TestWithContextAndFromContext(t *testing.T) {
	s := state.New("https://dns.example.com", "pwa_test")
	ctx := s.WithContext(context.Background())

	got := state.FromContext(ctx)
	if got == nil {
		t.Fatal("expected State from context, got nil")
	}
	if got.URL != s.URL {
		t.Errorf("URL = %q, want %q", got.URL, s.URL)
	}
}

func TestFromContextEmpty(t *testing.T) {
	got := state.FromContext(context.Background())
	if got == nil {
		t.Error("expected empty State, got nil")
	}
	if got.URL != "" {
		t.Errorf("expected empty URL, got: %q", got.URL)
	}
}
