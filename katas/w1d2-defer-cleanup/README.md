# W1D2: defer and resource cleanup

Timebox: 20 minutes for the fix, 5 for the prediction.

## Part A: predict, then run

Write down what this prints BEFORE you run it. Then run it in the
[Go playground](https://go.dev/play/) or a scratch `main.go`.

```go
func main() {
	x := 1
	defer fmt.Println("A:", x)
	x = 2
	defer func() { fmt.Println("B:", x) }()
	x = 3
	for i := 0; i < 3; i++ {
		defer fmt.Println("C:", i)
	}
}
```

## Part B: fix `SumSizes`

`SumSizes` opens every named resource and adds up their sizes.
It "works" on a small list, but it has bugs that would hurt in production.

Make every test pass without editing the tests:

```
go test -race ./w1d2-defer-cleanup/
```

## Explain it (3 to 5 sentences each)

1. When does a deferred call run, and in what order?
2. In Part A, why do A and B print different values for `x`?
3. Why is `defer` inside a loop dangerous here? What happens with 10,000 files?
4. Why should a `Close()` error on a file you wrote to never be ignored?
5. Transfer: what does Python's `with` or Java's try-with-resources give you that `defer` doesn't, and the other way round?

## Beyond Go (10 min)

Same bug, Python edition. Fix `beyond/cleanup.py`:

```
cd w1d2-defer-cleanup/beyond && python3 -m pytest -q
```

Then answer: Go needed a named return and a deferred closure to report the
close error. What did Python's `with` do for you there, for free?

## Rust, optional (15 min)

```
cd w1d2-defer-cleanup/rust && cargo test
```

Same bugs, fixed in `src/lib.rs`. Read the fake's `impl Drop` in the tests first.
Then answer:

1. Rust has no `defer`. What runs cleanup when a value goes out of scope, and
   what is that pattern called?
2. One test passed on the starter even though nothing called `close()` on the
   error path. Why?
3. `Drop::drop` can't return an error. So how do you report a close error in
   Rust? Compare that with Go's named return trick and Python's `with`.
