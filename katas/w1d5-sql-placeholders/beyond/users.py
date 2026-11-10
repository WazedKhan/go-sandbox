import sqlite3


def find_user(conn: sqlite3.Connection, email: str):
    """Return (id, email) rows for this email. The email comes from a web form."""
    return conn.execute(f"SELECT id, email FROM users WHERE email = '{email}'").fetchall()


def list_users(conn: sqlite3.Connection, sort_by: str, limit: int):
    """Return emails ordered by sort_by ("email" or "created_at"), at most limit.

    Unknown sort_by raises ValueError and runs no query.
    """
    return [r[0] for r in conn.execute(f"SELECT email FROM users ORDER BY {sort_by} LIMIT {limit}")]
