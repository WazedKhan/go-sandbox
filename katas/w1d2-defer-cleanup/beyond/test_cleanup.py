import pytest

from cleanup import sum_sizes


class FakeFS:
    def __init__(self, sizes, open_err=None, close_err=None):
        self.sizes = sizes
        self.open_err = open_err or {}
        self.close_err = close_err or {}
        self.open_now = 0
        self.max_open = 0

    def open(self, name):
        if name in self.open_err:
            raise self.open_err[name]
        self.open_now += 1
        self.max_open = max(self.max_open, self.open_now)
        return FakeRes(self, name)


class FakeRes:
    def __init__(self, fs, name):
        self.fs, self.name, self.closed = fs, name, False

    def size(self):
        return self.fs.sizes[self.name]

    def close(self):
        if not self.closed:
            self.closed = True
            self.fs.open_now -= 1
        if self.name in self.fs.close_err:
            raise self.fs.close_err[self.name]

    def __enter__(self):
        return self

    def __exit__(self, *exc):
        self.close()
        return False


def test_total():
    fs = FakeFS({"a": 10, "b": 20, "c": 30})
    assert sum_sizes(fs.open, ["a", "b", "c"]) == 60


def test_one_open_at_a_time():
    fs = FakeFS({n: 1 for n in "abcd"})
    sum_sizes(fs.open, list("abcd"))
    assert fs.max_open == 1
    assert fs.open_now == 0


def test_open_error_leaks_nothing():
    fs = FakeFS({"a": 1, "b": 1}, open_err={"c": PermissionError("c")})
    with pytest.raises(PermissionError):
        sum_sizes(fs.open, list("abcd"))
    assert fs.open_now == 0


def test_close_error_reaches_caller():
    fs = FakeFS({n: 1 for n in "abc"}, close_err={"b": OSError("disk full on flush")})
    with pytest.raises(OSError, match="disk full"):
        sum_sizes(fs.open, list("abc"))
    assert fs.open_now == 0
