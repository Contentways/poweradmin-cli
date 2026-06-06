package records_test

import (
	"testing"

	"github.com/contentways/poweradmin-cli/internal/cmd/records"
)

func TestNewRecordsCommand(t *testing.T) {
	cmd := records.NewRecordsCommand()

	if cmd.Use != "records" {
		t.Fatalf("Use = %q", cmd.Use)
	}

	if len(cmd.Commands()) != 3 {
		t.Fatalf("got %d subcommands, want 3", len(cmd.Commands()))
	}
}
