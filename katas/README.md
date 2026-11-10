# Katas

Small drills from The Engineer roadmap. One folder per roadmap day (`w1d2` = week 1, day 2).
Go is the tool, the idea is the lesson: every kata ends with a **Beyond Go** part
in Python, SQL, or plain HTTP. W1D1 to W1D3 also have an optional **Rust** part.

## Week 1

| Kata | Idea | Time |
| --- | --- | --- |
| `w1d1-map-aggregate` | Hash maps, zero values, nil maps | 10 min + 10 Python (+10 Rust) |
| `w1d2-defer-cleanup` | Resource cleanup, defer order, close errors | 25 min + 10 Python (+15 Rust) |
| `w1d3-waitgroup` | Goroutines, WaitGroup, data races | 25 min + 15 Python (+20 Rust) |
| `w1d4-context-timeout` | Deadlines and cancellation on a downstream call | 30 min + 15 Python |
| `w1d5-sql-placeholders` | Placeholders, injection, what can't be a placeholder | 30 min + 15 Python |
| `w1-weekend-ai-review` | Review AI-written code like a senior | 90 min |

## How to work a kata

1. Read the kata README. No AI, no docs for the first attempt.
2. Run the tests and read the failures: `go test -race ./w1d2-defer-cleanup/`
3. Fix the code until they pass. Don't edit the tests.
4. Answer the "Explain it" questions in your own words and send them to Claude.
5. Do the Beyond Go part: `cd <kata>/beyond && python3 -m pytest -q`
   Optional Rust: `cd <kata>/rust && cargo test`
6. Commit your solution as its own commit.

Run from this `katas/` folder; it is its own Go module.
