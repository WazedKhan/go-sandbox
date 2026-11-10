# W1D5: SQL placeholders and injection

Timebox: 30 minutes.

## Task

Fix both functions in `users.go`. The tests use a fake driver
(`fakedb_test.go`) that records the exact SQL text and arguments sent,
so they can see whether user input ends up inside the SQL.

```
go test -race ./w1d5-sql-placeholders/
```

Placeholders are Postgres style: `$1`, `$2`, ...

There's one trap the interview didn't cover: some parts of a query
**can't** be placeholders. Find out which, and how to handle them safely.

## Explain it (3 to 5 sentences each)

1. Why does a placeholder stop injection, while escaping quotes yourself doesn't?
   (Think: who parses the SQL, and when does your value arrive?)
2. Why can't `ORDER BY $1` work? What did you do instead?
3. `ListUsers` had two bugs that had nothing to do with injection. What were
   they, and which earlier kata do they connect to?
4. Is `fmt.Sprintf` with `%d` for an int "safe enough"? Argue both sides.

## Beyond Go (15 min): a real injection

`beyond/users.py` uses Python's built-in `sqlite3` with a real database.
First run the tests and watch the attack actually dump every user. Then fix it.

```
cd w1d5-sql-placeholders/beyond && python3 -m pytest -q
```

Then answer: SQLite uses `?`, Postgres uses `$1`, Django uses `%s`. What is
the same underneath all three? And where does Django's ORM still let you
inject (look up `.extra()` and `.raw()`)?
