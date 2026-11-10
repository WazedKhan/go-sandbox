# W1D3: goroutines, WaitGroup, race detector

Timebox: 25 minutes.

## Task

Fix both functions in `fanout.go`. The tests run with `-race`, so a
data race is a failure even when the numbers look right.

```
go test -race ./w1d3-waitgroup/
```

1. `FetchAll` must call `fetch` for every id **concurrently**, wait for all of
   them, and return results **in the same order as ids**.
2. `CountParallel` starts `n` goroutines that each add 1 to a shared counter.
   It must return exactly `n`, every time.

First, run the tests on the starter code once and read what the race
detector prints. That output is the lesson.

Stretch: your Go is 1.26, so after you pass, try rewriting `FetchAll` with
`wg.Go(func() { ... })` (added in Go 1.25). Compare it to `Add`/`Done`.
If `go test` complains about the version, change the `go` line in
`katas/go.mod` to `go 1.25`.

## Explain it (3 to 5 sentences each)

1. Why did the starter `FetchAll` often return an empty or short slice?
2. Why must `wg.Add` happen before `go func()`, not inside it?
3. Writing to `results[i]` from many goroutines is safe, but `append` is not. Why?
4. Mutex vs `sync/atomic` for the counter: when would you pick each?
5. Transfer: goroutines vs OS threads vs Python `asyncio`. What is each one good at?

## Beyond Go (15 min)

Make `beyond/fanout.py` concurrent with the standard library only:

```
cd w1d3-waitgroup/beyond && python3 -m pytest -q
```

Then answer: Python threads share a GIL. Why did threads still make this
20x faster? What kind of work would they NOT speed up, and what would you use then?

## Rust, optional (20 min)

```
cd w1d3-waitgroup/rust && cargo test
```

Before fixing anything, write `count_parallel` the way the Go starter did:
a plain `let mut count = 0;` that every thread increments. Try to compile it
and read the error. Then fix both functions, standard library only.

1. Go found the race at run time with `-race`. When did Rust find it?
2. What do `Send` and `Sync` promise, in one sentence each?
3. `std::thread::scope` plays which role from Go's `sync` package?
4. Why are OS threads fine for 20 calls but a bad idea for 100,000?
