package titles

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

type bodyTracker struct {
	next   http.RoundTripper
	opened atomic.Int64
	closed atomic.Int64
}

type trackedBody struct {
	io.ReadCloser
	t    *bodyTracker
	once atomic.Bool
}

func (b *trackedBody) Close() error {
	if b.once.CompareAndSwap(false, true) {
		b.t.closed.Add(1)
	}
	return b.ReadCloser.Close()
}

func (t *bodyTracker) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := t.next.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	t.opened.Add(1)
	resp.Body = &trackedBody{ReadCloser: resp.Body, t: t}
	return resp, nil
}

func newClient() (*http.Client, *bodyTracker) {
	t := &bodyTracker{next: http.DefaultTransport}
	return &http.Client{Transport: t}, t
}

func server() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/broken":
			http.Error(w, "boom", http.StatusInternalServerError)
		case "/slow":
			select {
			case <-r.Context().Done():
			case <-time.After(3 * time.Second):
			}
		default:
			time.Sleep(20 * time.Millisecond)
			fmt.Fprintf(w, "  Title of %s \n", r.URL.Path)
		}
	}))
}

func TestFetchTitles(t *testing.T) {
	srv := server()
	defer srv.Close()
	client, tr := newClient()

	var urls []string
	want := map[string]string{}
	for i := 0; i < 30; i++ {
		u := fmt.Sprintf("%s/page-%d", srv.URL, i)
		urls = append(urls, u)
		want[u] = fmt.Sprintf("Title of /page-%d", i)
	}

	got, err := FetchTitles(context.Background(), client, urls)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !maps.Equal(got, want) {
		t.Errorf("got %d titles, want %d (or values differ)", len(got), len(want))
	}
	if tr.opened.Load() != int64(len(urls)) {
		t.Errorf("only %d of %d requests went through the given client", tr.opened.Load(), len(urls))
	}
	if tr.closed.Load() != tr.opened.Load() {
		t.Errorf("leaked bodies: opened %d, closed %d", tr.opened.Load(), tr.closed.Load())
	}
}

func TestFetchTitlesHTTPErrorIsAnError(t *testing.T) {
	srv := server()
	defer srv.Close()
	client, _ := newClient()

	urls := []string{srv.URL + "/a", srv.URL + "/broken", srv.URL + "/b"}
	_, err := FetchTitles(context.Background(), client, urls)
	if err == nil {
		t.Fatal(`a 500 response came back as a "title"; it should be an error`)
	}
}

func TestFetchTitlesRespectsContext(t *testing.T) {
	srv := server()
	defer srv.Close()
	client, _ := newClient()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := FetchTitles(ctx, client, []string{srv.URL + "/a", srv.URL + "/slow"})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want context.DeadlineExceeded", err)
	}
	if time.Since(start) > time.Second {
		t.Errorf("took %v; ctx was ignored", time.Since(start))
	}
}
