import sqlite3

import pytest

from users import find_user, list_users


@pytest.fixture
def conn():
    c = sqlite3.connect(":memory:")
    c.execute("CREATE TABLE users (id INTEGER PRIMARY KEY, email TEXT, password_hash TEXT, created_at INTEGER)")
    c.executemany(
        "INSERT INTO users (email, password_hash, created_at) VALUES (?, ?, ?)",
        [("alice@example.com", "h1", 3), ("bob@example.com", "h2", 1), ("wazed@example.com", "h3", 2)],
    )
    yield c
    c.close()


def test_find_user(conn):
    assert find_user(conn, "bob@example.com") == [(2, "bob@example.com")]


def test_find_user_injection_dumps_nothing(conn):
    assert find_user(conn, "x' OR '1'='1") == []


def test_find_user_quote_in_email_is_fine(conn):
    # a real email can contain an apostrophe; it must not break the query
    assert find_user(conn, "o'brien@example.com") == []


def test_list_users_sorted(conn):
    assert list_users(conn, "created_at", 2) == ["bob@example.com", "wazed@example.com"]


@pytest.mark.parametrize("bad", ["password_hash", "email; DROP TABLE users", ""])
def test_list_users_rejects_unknown_sort(conn, bad):
    with pytest.raises(ValueError):
        list_users(conn, bad, 10)
    assert conn.execute("SELECT count(*) FROM users").fetchone()[0] == 3


def test_list_users_limit_is_not_sql(conn):
    with pytest.raises((ValueError, TypeError, sqlite3.Error)):
        list_users(conn, "email", "1 UNION SELECT password_hash FROM users")
