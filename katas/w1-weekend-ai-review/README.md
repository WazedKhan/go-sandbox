# Week 1 weekend: review the AI's code

Timebox: 90 minutes. This is the core skill of AI-assisted engineering:
AI writes code fast, and **you** decide whether it is safe to ship.

`titles.go` was written by an AI. It compiles and "works" on a quick try.

## Step 1: read only (20 min, no running, no AI)

Write `REVIEW.md` in this folder as PR review comments, one per problem:

```
L27: <what's wrong> -> <what breaks in production> -> <the fix>
```

Aim for at least 5. Every week 1 kata is hiding in this file.

## Step 2: run the tests

```
go test -race ./w1-weekend-ai-review/
```

Which problems did the tests catch that you missed? Which did you catch that
the tests don't check? (A good reviewer finds things tests never will.)

## Step 3: fix it yourself

Make the tests pass with `-race`. Keep it concurrent.
Stretch: cap concurrency at 5 requests in flight. What would you use?

## Step 4: AI vs you (15 min)

Paste the original `titles.go` into any AI and ask it to "fix the bugs".
Review its fix against your `REVIEW.md`. Did it catch everything? Did it
introduce anything new? Write 3 sentences on what you'd trust it with
and what you wouldn't.

## Beyond Go

Every problem in this file has a twin in Python, Java, or Node. For 3 of
your review comments, add one line each: "Same bug in Python: ..."
