package transform

import (
	"context"
	"path/filepath"
	"runtime/trace"
	"testing"

	"github.com/rusq/fsadapter"
	"github.com/rusq/slackdump/v2/internal/chunk/state"
	"github.com/rusq/slackdump/v2/types"
)

const whereTheTempIsAt = "../../tmp"

func TestStandard_Transform(t *testing.T) {
	ctx, task := trace.NewTask(context.Background(), "TestStandard_Transform")
	defer task.End()
	// MANUAL
	fs := fsadapter.NewDirectory(filepath.Join(whereTheTempIsAt, "manual"))
	s := NewStandard(fs, func(c *types.Conversation) string {
		return c.ID + ".json"
	})
	st, err := state.Load(filepath.Join(whereTheTempIsAt, "C01SPFM1KNY.state"))
	if err != nil {
		t.Fatalf("state.Load(): %s", err)
	}
	if err := s.Transform(ctx, st, whereTheTempIsAt); err != nil {
		t.Fatal(err)
	}
}
