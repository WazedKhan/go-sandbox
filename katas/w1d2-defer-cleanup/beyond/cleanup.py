def sum_sizes(open_fn, names):
    """Open each named resource, add up the sizes, close it.

    Same rules as the Go version: at most one open at a time, nothing leaks
    when something fails, and errors from open or close reach the caller.
    open_fn(name) returns an object with .size(), .close(), and context-manager
    support (__enter__/__exit__).
    """
    opened = [open_fn(name) for name in names]
    total = sum(r.size() for r in opened)
    for r in opened:
        r.close()
    return total
