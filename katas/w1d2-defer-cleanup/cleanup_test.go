package cleanup

import (
	"errors"
	"strings"
	"testing"
)

// fakeFS tracks how many resources are open at once.
type fakeFS struct {
	sizes    map[string]int
	openErr  map[string]error
	closeErr map[string]error
	open     int
	maxOpen  int
	closed   []string
}

type fakeRes struct {
	fs   *fakeFS
	name string
}

func (r *fakeRes) Size() int { return r.fs.sizes[r.name] }

func (r *fakeRes) Close() error {
	r.fs.open--
	r.fs.closed = append(r.fs.closed, r.name)
	return r.fs.closeErr[r.name]
}

func (fs *fakeFS) Open(name string) (Resource, error) {
	if err := fs.openErr[name]; err != nil {
		return nil, err
	}
	fs.open++
	if fs.open > fs.maxOpen {
		fs.maxOpen = fs.open
	}
	return &fakeRes{fs: fs, name: name}, nil
}

func TestSumSizesTotal(t *testing.T) {
	fs := &fakeFS{sizes: map[string]int{"a": 10, "b": 20, "c": 30}}
	got, err := SumSizes(fs.Open, []string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 60 {
		t.Errorf("total = %d, want 60", got)
	}
}

func TestSumSizesOneOpenAtATime(t *testing.T) {
	fs := &fakeFS{sizes: map[string]int{"a": 1, "b": 1, "c": 1, "d": 1}}
	if _, err := SumSizes(fs.Open, []string{"a", "b", "c", "d"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fs.maxOpen != 1 {
		t.Errorf("max resources open at once = %d, want 1", fs.maxOpen)
	}
	if fs.open != 0 {
		t.Errorf("%d resources still open after return", fs.open)
	}
}

func TestSumSizesOpenErrorClosesEarlier(t *testing.T) {
	boom := errors.New("permission denied")
	fs := &fakeFS{
		sizes:   map[string]int{"a": 1, "b": 1},
		openErr: map[string]error{"c": boom},
	}
	_, err := SumSizes(fs.Open, []string{"a", "b", "c", "d"})
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want it to wrap %v", err, boom)
	}
	if !strings.Contains(err.Error(), "c") {
		t.Errorf("error %q should name the resource that failed", err)
	}
	if fs.open != 0 {
		t.Errorf("%d resources leaked after an open error", fs.open)
	}
}

func TestSumSizesCloseErrorIsReturned(t *testing.T) {
	flush := errors.New("disk full on flush")
	fs := &fakeFS{
		sizes:    map[string]int{"a": 1, "b": 1, "c": 1},
		closeErr: map[string]error{"b": flush},
	}
	_, err := SumSizes(fs.Open, []string{"a", "b", "c"})
	if !errors.Is(err, flush) {
		t.Fatalf("err = %v, want it to wrap the close error %v", err, flush)
	}
	if !strings.Contains(err.Error(), "b") {
		t.Errorf("error %q should name the resource that failed", err)
	}
	if fs.open != 0 {
		t.Errorf("%d resources leaked after a close error", fs.open)
	}
}
