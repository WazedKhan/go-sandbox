package users

// A tiny database/sql driver that records every query and its arguments
// instead of talking to Postgres. Not part of the kata; don't edit it.

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"sync"
)

type call struct {
	query string
	args  []driver.Value
}

type fakeDB struct {
	mu         sync.Mutex
	calls      []call
	rows       []any // values for "SELECT name" queries; an error value fails Next
	openRows   int
	closedRows int
}

func (f *fakeDB) open() *sql.DB { return sql.OpenDB(f) }

func (f *fakeDB) Connect(context.Context) (driver.Conn, error) { return &fakeConn{f}, nil }
func (f *fakeDB) Driver() driver.Driver                        { return nil }

type fakeConn struct{ f *fakeDB }

func (c *fakeConn) Prepare(q string) (driver.Stmt, error) { return &fakeStmt{c.f, q}, nil }
func (c *fakeConn) Close() error                          { return nil }
func (c *fakeConn) Begin() (driver.Tx, error)             { return nil, errors.New("no tx") }

type fakeStmt struct {
	f *fakeDB
	q string
}

func (s *fakeStmt) Close() error  { return nil }
func (s *fakeStmt) NumInput() int { return -1 }
func (s *fakeStmt) Exec([]driver.Value) (driver.Result, error) {
	return nil, errors.New("exec not supported")
}

func (s *fakeStmt) Query(args []driver.Value) (driver.Rows, error) {
	s.f.mu.Lock()
	defer s.f.mu.Unlock()
	s.f.calls = append(s.f.calls, call{s.q, args})
	s.f.openRows++
	var data []any
	if len(s.q) >= 9 && s.q[:9] == "SELECT id" {
		data = []any{int64(42)}
	} else {
		data = s.f.rows
	}
	return &fakeRows{f: s.f, data: data}, nil
}

type fakeRows struct {
	f      *fakeDB
	data   []any
	i      int
	closed bool
}

func (r *fakeRows) Columns() []string { return []string{"col"} }

func (r *fakeRows) Close() error {
	r.f.mu.Lock()
	defer r.f.mu.Unlock()
	if !r.closed {
		r.closed = true
		r.f.closedRows++
	}
	return nil
}

func (r *fakeRows) Next(dest []driver.Value) error {
	if r.i >= len(r.data) {
		return io.EOF
	}
	v := r.data[r.i]
	r.i++
	if err, ok := v.(error); ok {
		return err
	}
	dest[0] = v
	return nil
}
