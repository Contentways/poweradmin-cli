package users_test

import (
	"testing"

	"github.com/contentways/poweradmin-cli/internal/cmd/users"
)

func TestNewUsersCommand(t *testing.T) {
	cmd := users.NewUsersCommand()

	if cmd.Use != "users" {
		t.Fatalf("Use = %q", cmd.Use)
	}

	if len(cmd.Commands()) != 6 {
		t.Fatalf("got %d subcommands, want 6", len(cmd.Commands()))
	}
}
