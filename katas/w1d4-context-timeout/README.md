# W1D4: context deadlines and cancellation

Timebox: 30 minutes.

This is the "downstream API hangs" interview question, as code.

## Task

Fix `FetchProfile` in `profile.go`. It calls `GET {baseURL}/users/{id}`
on a downstream service. Rules are in the doc comment.

```
go test -race ./w1d4-context-timeout/
```

The starter passes the happy path, which is exactly why it's dangerous:
it looks fine until the downstream gets slow.

## Explain it (3 to 5 sentences each)

1. Why is `context` better than a bare `time.After` or `time.Sleep` timer here?
2. The caller already passes a `ctx`. Why add your own `CallTimeout` on top?
   What happens when both exist?
3. What goes wrong in a busy service if you forget `resp.Body.Close()`?
4. Why use the `*http.Client` passed in instead of `http.Get`?
5. Your real service: write 5 sentences on what happens today when one of
   your legacy downstream APIs hangs for 30 seconds. Who waits, what fills up,
   what the user sees. (This is the note you still owe from Day 1.)

## Beyond Go (15 min): the timeout that isn't

Read `beyond/slow_download.py`. It calls a slow downstream with `timeout=1`.
**Before running it, write down how long you think the call takes.** Then:

```
cd w1d4-context-timeout/beyond && python3 slow_download.py
```

Answer:

1. Why didn't `timeout=1` stop the call at 1 second? (`requests` behaves the same way.)
2. What does that mean for a Django view calling a slow downstream with `requests`?
3. Go's `context.WithTimeout` didn't have this problem. Why?
4. How would you get a real total deadline in Python? Look it up, then answer
   in 3 sentences.
