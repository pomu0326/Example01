package todo

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func newTestStore(t *testing.T) *SQLiteStore {
	t.Helper()
	s, err := OpenSQLite(":memory:")
	if err != nil {
		t.Fatalf("OpenSQLite: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	// Deterministic, strictly increasing clock so ordering is stable.
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	var tick int
	s.now = func() time.Time { tick++; return base.Add(time.Duration(tick) * time.Second) }
	return s
}

func ptr[T any](v T) *T { return &v }

func TestCreateAndGet(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	created, err := s.Create(ctx, "  buy milk  ")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == 0 || created.Title != "buy milk" || created.Completed {
		t.Fatalf("unexpected todo: %+v", created)
	}

	got, err := s.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != created {
		t.Fatalf("Get = %+v, want %+v", got, created)
	}
}

func TestCreateValidation(t *testing.T) {
	s := newTestStore(t)
	for _, title := range []string{"", "   ", strings.Repeat("あ", MaxTitleLength+1)} {
		if _, err := s.Create(context.Background(), title); !errors.Is(err, ErrInvalidTitle) {
			t.Errorf("Create(%q) err = %v, want ErrInvalidTitle", title, err)
		}
	}
	if _, err := s.Create(context.Background(), strings.Repeat("あ", MaxTitleLength)); err != nil {
		t.Errorf("Create(max length) err = %v", err)
	}
}

func TestListNewestFirst(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	empty, err := s.List(ctx)
	if err != nil || len(empty) != 0 {
		t.Fatalf("List on empty store = %v, %v", empty, err)
	}

	for _, title := range []string{"a", "b", "c"} {
		if _, err := s.Create(ctx, title); err != nil {
			t.Fatal(err)
		}
	}
	todos, err := s.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var titles []string
	for _, td := range todos {
		titles = append(titles, td.Title)
	}
	if strings.Join(titles, ",") != "c,b,a" {
		t.Fatalf("List order = %v, want [c b a]", titles)
	}
}

func TestUpdate(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	created, _ := s.Create(ctx, "old")

	updated, err := s.Update(ctx, created.ID, Patch{Completed: ptr(true)})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if !updated.Completed || updated.Title != "old" {
		t.Fatalf("completed-only update changed wrong fields: %+v", updated)
	}
	if !updated.UpdatedAt.After(created.UpdatedAt) {
		t.Fatalf("UpdatedAt not advanced: %v -> %v", created.UpdatedAt, updated.UpdatedAt)
	}

	updated, err = s.Update(ctx, created.ID, Patch{Title: ptr(" new ")})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Title != "new" || !updated.Completed {
		t.Fatalf("title-only update changed wrong fields: %+v", updated)
	}

	if _, err := s.Update(ctx, created.ID, Patch{Title: ptr("")}); !errors.Is(err, ErrInvalidTitle) {
		t.Fatalf("Update with empty title err = %v", err)
	}
	if _, err := s.Update(ctx, 999, Patch{Completed: ptr(true)}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Update missing err = %v", err)
	}
}

func TestDelete(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	created, _ := s.Create(ctx, "x")

	if err := s.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Get(ctx, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get after delete err = %v", err)
	}
	if err := s.Delete(ctx, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("second Delete err = %v", err)
	}
}
