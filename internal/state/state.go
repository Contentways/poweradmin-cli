// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package state

import (
	"context"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
)

// State holds the global runtime state passed to all commands.
type State struct {
	client *poweradmin.Client
}

// New creates a new State with a configured Poweradmin client.
func New(url, apiKey string) (*State, error) {
	client, err := poweradmin.NewClient(
		poweradmin.WithBaseURL(url),
		poweradmin.WithAPIKey(apiKey),
	)
	if err != nil {
		return nil, err
	}
	return &State{client: client}, nil
}

// Client returns the Poweradmin client.
func (s *State) Client() *poweradmin.Client {
	return s.client
}

// WithContext returns a context with the state attached.
func (s *State) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, stateKey{}, s)
}

type stateKey struct{}

// FromContext retrieves the State from a context.
func FromContext(ctx context.Context) *State {
	s, _ := ctx.Value(stateKey{}).(*State)
	return s
}
