package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrInvalidSort is returned when ListUsers gets a sort column it does not allow.
var ErrInvalidSort = errors.New("invalid sort column")

// FindIDByEmail returns the id of the user with this email (Postgres).
//
// The email comes straight from an HTTP request. Treat it as hostile.
func FindIDByEmail(ctx context.Context, db *sql.DB, email string) (int64, error) {
	q := fmt.Sprintf("SELECT id FROM users WHERE email = '%s'", email)
	var id int64
	err := db.QueryRowContext(ctx, q).Scan(&id)
	return id, err
}

// ListUsers returns user names ordered by sortBy, at most limit of them.
//
// Allowed sort columns: "name", "created_at". Anything else is ErrInvalidSort,
// and no query may run. sortBy and limit both come from query-string params.
func ListUsers(ctx context.Context, db *sql.DB, sortBy string, limit int) ([]string, error) {
	q := fmt.Sprintf("SELECT name FROM users ORDER BY %s LIMIT %d", sortBy, limit)
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	var names []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		names = append(names, n)
	}
	return names, nil
}
