import time

from fanout import fetch_all


def test_order():
    def fetch(i):
        time.sleep((10 - i) / 1000)  # later ids finish first
        return f"user-{i}"

    assert fetch_all([5, 1, 4, 2, 3], fetch) == ["user-5", "user-1", "user-4", "user-2", "user-3"]


def test_concurrent():
    def fetch(i):
        time.sleep(0.1)  # I/O wait, like an HTTP call
        return "ok"

    start = time.monotonic()
    got = fetch_all(list(range(20)), fetch)
    elapsed = time.monotonic() - start
    assert len(got) == 20
    assert elapsed < 0.6, f"took {elapsed:.2f}s; calls are not concurrent"


def test_empty():
    assert fetch_all([], lambda i: i) == []
