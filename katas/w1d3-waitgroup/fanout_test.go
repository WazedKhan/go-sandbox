package fanout

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

func TestFetchAllOrder(t *testing.T) {
	ids := []int{5, 1, 4, 2, 3}
	fetch := func(id int) string {
		// later ids finish first, so order must come from the index, not timing
		time.Sleep(time.Duration(10-id) * time.Millisecond)
		return fmt.Sprintf("user-%d", id)
	}
	got := FetchAll(ids, fetch)
	want := []string{"user-5", "user-1", "user-4", "user-2", "user-3"}
	if !slices.Equal(got, want) {
		t.Errorf("FetchAll = %v, want %v", got, want)
	}
}

func TestFetchAllIsConcurrent(t *testing.T) {
	ids := make([]int, 20)
	for i := range ids {
		ids[i] = i
	}
	fetch := func(id int) string {
		time.Sleep(100 * time.Millisecond)
		return "ok"
	}
	start := time.Now()
	got := FetchAll(ids, fetch)
	elapsed := time.Since(start)

	if len(got) != len(ids) {
		t.Fatalf("got %d results, want %d", len(got), len(ids))
	}
	// 20 calls x 100ms = 2s if sequential. Concurrent should be ~100ms.
	if elapsed > 600*time.Millisecond {
		t.Errorf("took %v; calls are not running concurrently", elapsed)
	}
}

func TestFetchAllEmpty(t *testing.T) {
	got := FetchAll(nil, func(int) string { return "x" })
	if len(got) != 0 {
		t.Errorf("FetchAll(nil) = %v, want empty", got)
	}
}

func TestCountParallel(t *testing.T) {
	for _, n := range []int{1, 10, 1000} {
		if got := CountParallel(n); got != n {
			t.Errorf("CountParallel(%d) = %d", n, got)
		}
	}
}
