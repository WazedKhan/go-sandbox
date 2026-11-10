package users

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
)

const attack = "x' OR '1'='1"

func TestFindIDByEmailUsesPlaceholder(t *testing.T) {
	f := &fakeDB{}
	db := f.open()
	defer db.Close()

	id, err := FindIDByEmail(context.Background(), db, "wazed@example.com")
	if err != nil || id != 42 {
		t.Fatalf("got (%d, %v), want (42, nil)", id, err)
	}
	c := f.calls[0]
	if !strings.Contains(c.query, "$1") {
		t.Errorf("query %q has no $1 placeholder", c.query)
	}
	if strings.Contains(c.query, "wazed@example.com") {
		t.Errorf("user input was pasted into the SQL text: %q", c.query)
	}
	if len(c.args) != 1 || c.args[0] != "wazed@example.com" {
		t.Errorf("args = %v, want [wazed@example.com]", c.args)
	}
}

func TestFindIDByEmailInjection(t *testing.T) {
	f := &fakeDB{}
	db := f.open()
	defer db.Close()

	FindIDByEmail(context.Background(), db, "a@example.com")
	FindIDByEmail(context.Background(), db, attack)

	if f.calls[0].query != f.calls[1].query {
		t.Errorf("SQL text changes with user input:\n  %q\n  %q", f.calls[0].query, f.calls[1].query)
	}
	if strings.Contains(f.calls[1].query, "OR '1'='1") {
		t.Errorf("injection reached the SQL text: %q", f.calls[1].query)
	}
}

func TestListUsersAllowedSort(t *testing.T) {
	f := &fakeDB{rows: []any{"alice", "bob"}}
	db := f.open()
	defer db.Close()

	got, err := ListUsers(context.Background(), db, "created_at", 37)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !slices.Equal(got, []string{"alice", "bob"}) {
		t.Errorf("got %v", got)
	}
	c := f.calls[0]
	if !strings.Contains(c.query, "ORDER BY created_at") {
		t.Errorf("query %q should order by created_at", c.query)
	}
	if strings.Contains(c.query, "37") {
		t.Errorf("limit was pasted into the SQL text: %q", c.query)
	}
	if !slices.Contains(c.args, any(int64(37))) {
		t.Errorf("args = %v, want the limit 37 passed as a placeholder arg", c.args)
	}
}

func TestListUsersRejectsUnknownSort(t *testing.T) {
	for _, sortBy := range []string{"password", "name; DROP TABLE users", "", "NAME"} {
		f := &fakeDB{}
		db := f.open()

		_, err := ListUsers(context.Background(), db, sortBy, 10)
		if !errors.Is(err, ErrInvalidSort) {
			t.Errorf("sortBy %q: err = %v, want ErrInvalidSort", sortBy, err)
		}
		if len(f.calls) != 0 {
			t.Errorf("sortBy %q: a query ran anyway: %q", sortBy, f.calls[0].query)
		}
		db.Close()
	}
}

func TestListUsersClosesRowsOnScanError(t *testing.T) {
	f := &fakeDB{rows: []any{"alice", nil, "bob"}} // NULL can't scan into string
	db := f.open()
	defer db.Close()

	if _, err := ListUsers(context.Background(), db, "name", 10); err == nil {
		t.Fatal("expected a scan error for a NULL name")
	}
	if f.openRows != f.closedRows {
		t.Errorf("rows leaked: opened %d, closed %d", f.openRows, f.closedRows)
	}
}

func TestListUsersReportsIterationError(t *testing.T) {
	lost := errors.New("connection reset mid-result")
	f := &fakeDB{rows: []any{"alice", lost}}
	db := f.open()
	defer db.Close()

	got, err := ListUsers(context.Background(), db, "name", 10)
	if !errors.Is(err, lost) {
		t.Errorf("got (%v, %v); a broken result set must not look like success", got, err)
	}
}
