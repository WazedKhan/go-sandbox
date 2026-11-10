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
