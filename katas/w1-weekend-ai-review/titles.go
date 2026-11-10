package titles

// This code was written by an AI assistant from the prompt:
// "Write a Go function that fetches a list of URLs concurrently and returns
// a map of URL to page title. Handle errors."
//
// It compiles. It even works on a quick manual test. Review it like a
// senior engineer before anything runs.

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
)

// FetchTitles fetches every URL concurrently and returns url -> title,
// where the title is the response body with surrounding whitespace trimmed.
// If any fetch fails, it returns an error.
func FetchTitles(ctx context.Context, client *http.Client, urls []string) (map[string]string, error) {
	titles := make(map[string]string)
	var wg sync.WaitGroup
	var err error

	for _, url := range urls {
		go func() {
			wg.Add(1)
			defer wg.Done()

			resp, fetchErr := http.Get(url)
			if fetchErr != nil {
				err = fetchErr
				return
			}
			body, _ := io.ReadAll(resp.Body)
			titles[url] = strings.TrimSpace(string(body))
		}()
	}

	wg.Wait()
	return titles, err
}
