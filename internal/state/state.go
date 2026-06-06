// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package state

import (
	"context"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
)

type State struct {
	URL    string
	APIKey string
}

func New(url, apiKey string) *State {
	return &State{
		URL:    url,
		APIKey: apiKey,
	}
}

func (s *State) Client() (*poweradmin.Client, error) {
	return poweradmin.NewClient(
		poweradmin.WithBaseURL(s.URL),
		poweradmin.WithAPIKey(s.APIKey),
	)
}

func (s *State) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, stateKey{}, s)
}

type stateKey struct{}

func FromContext(ctx context.Context) *State {
	s, _ := ctx.Value(stateKey{}).(*State)
	return s
}
