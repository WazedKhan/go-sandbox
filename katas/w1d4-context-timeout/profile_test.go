package profile

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// bodyTracker counts response bodies opened and closed through the client.
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

func okServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/7" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"id":7,"name":"Wazed"}`)
	}))
}

// hangingServer never answers until the client gives up (or 3s pass).
func hangingServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(3 * time.Second):
			io.WriteString(w, `{"id":7,"name":"too late"}`)
		}
	}))
}

func TestFetchProfileOK(t *testing.T) {
	srv := okServer()
	defer srv.Close()
	client, tr := newClient()

	p, err := FetchProfile(context.Background(), client, srv.URL, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != (Profile{ID: 7, Name: "Wazed"}) {
		t.Errorf("got %+v", p)
	}
	if tr.opened.Load() != 1 {
		t.Fatalf("the passed-in client was not used (%d requests through it)", tr.opened.Load())
	}
	if tr.closed.Load() != 1 {
		t.Errorf("response body was not closed")
	}
}

func TestFetchProfileNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down for maintenance", http.StatusServiceUnavailable)
	}))
	defer srv.Close()
	client, tr := newClient()

	_, err := FetchProfile(context.Background(), client, srv.URL, 7)
	if err == nil {
		t.Fatal("expected an error for a 503 response")
	}
	if !strings.Contains(err.Error(), "503") {
		t.Errorf("error %q should include the status code", err)
	}
	if tr.closed.Load() != tr.opened.Load() {
		t.Errorf("response body was not closed on the error path")
	}
}

func TestFetchProfileCallerDeadline(t *testing.T) {
	srv := hangingServer()
	defer srv.Close()
	client, _ := newClient()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := FetchProfile(ctx, client, srv.URL, 7)
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want context.DeadlineExceeded", err)
	}
	if elapsed > 300*time.Millisecond {
		t.Errorf("took %v; the caller's 50ms deadline was ignored", elapsed)
	}
}

func TestFetchProfileOwnTimeout(t *testing.T) {
	srv := hangingServer()
	defer srv.Close()
	client, _ := newClient()

	start := time.Now()
	_, err := FetchProfile(context.Background(), client, srv.URL, 7)
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want context.DeadlineExceeded", err)
	}
	if elapsed > CallTimeout+400*time.Millisecond {
		t.Errorf("took %v; CallTimeout (%v) was not enforced", elapsed, CallTimeout)
	}
}

func TestFetchProfileCanceled(t *testing.T) {
	srv := hangingServer()
	defer srv.Close()
	client, _ := newClient()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel() // e.g. the user closed the browser tab
	}()

	start := time.Now()
	_, err := FetchProfile(ctx, client, srv.URL, 7)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
	if time.Since(start) > 300*time.Millisecond {
		t.Errorf("took %v; cancellation was ignored", time.Since(start))
	}
}
