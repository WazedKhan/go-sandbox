package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Profile is what the downstream user service returns.
type Profile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CallTimeout is the longest one downstream call may take, even if the
// caller's ctx has no deadline at all.
const CallTimeout = 500 * time.Millisecond

// FetchProfile calls GET {baseURL}/users/{id} with the given client.
//
// Rules:
//   - Stop waiting as soon as ctx is done (deadline or cancel).
//   - Never wait longer than CallTimeout, even when ctx has no deadline.
//   - The returned error must still match context.DeadlineExceeded or
//     context.Canceled with errors.Is when that is the cause.
//   - A non-200 response is an error that includes the status code.
//   - Always close the response body.
func FetchProfile(ctx context.Context, client *http.Client, baseURL string, id int) (Profile, error) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%d", baseURL, id))
	if err != nil {
		return Profile{}, err
	}

	var p Profile
	err = json.NewDecoder(resp.Body).Decode(&p)
	return p, err
}
