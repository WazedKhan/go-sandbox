# W1D1: Map aggregation

Timebox: 10 minutes.

## Task

Implement `Aggregate` in `aggregate.go`.

```
input:  [{"m1": 20}, {"m2": 30}, {"m1": 20}]
output: {"m1": 40, "m2": 30}
```

Make every test pass:

```
go test -race ./w1d1-map-aggregate/
```

## Explain it (write 3 to 5 sentences, then send them to Claude)

1. What does `counts["missing"]` return, and why does that make your loop short?
2. Reading from a nil map is fine. Writing to one panics. Why the difference?
3. Why does `range` over a map give a different order each run?
4. Transfer: how would you write this in Python? What plays the role of Go's zero value there?

## Beyond Go (10 min)

Same problem in Python. Implement `beyond/aggregate.py`:

```
cd w1d1-map-aggregate/beyond && python3 -m pytest -q
```

Try it with `collections.Counter` and `+` first, and watch which test fails.
Then explain in one sentence why `test_returns_plain_dict` exists.

## Rust, optional (10 min)

```
cd w1d1-map-aggregate/rust && cargo test
```

Implement `aggregate` in `src/lib.rs`, standard library only. Then answer:

1. Go has a zero value and Python has `defaultdict`. What does Rust use for
   "add to this key, starting from 0"?
2. Rust has no nil map. What does `map.get("missing")` return, and why does
   that rule out Go's nil-map panic?
3. Why does `aggregate` take `&[HashMap<..>]` and not `Vec<HashMap<..>>`?
   What would the caller lose with `Vec`?
