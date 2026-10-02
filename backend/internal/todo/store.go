package todo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "modernc.org/sqlite" // registers the "sqlite" driver
)

// Store persists todos.
type Store interface {
	List(ctx context.Context) ([]Todo, error)
	Get(ctx context.Context, id int64) (Todo, error)
	Create(ctx context.Context, title string) (Todo, error)
	Update(ctx context.Context, id int64, p Patch) (Todo, error)
	Delete(ctx context.Context, id int64) error
}

const schema = `
CREATE TABLE IF NOT EXISTS todos (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	title      TEXT    NOT NULL,
	completed  INTEGER NOT NULL DEFAULT 0,
	created_at TEXT    NOT NULL,
	updated_at TEXT    NOT NULL
)`

// timeLayout is fixed-width so that timestamps sort correctly as strings.
const timeLayout = "2006-01-02T15:04:05.000000000Z07:00"

// SQLiteStore is a Store backed by SQLite.
type SQLiteStore struct {
	db  *sql.DB
	now func() time.Time
}

// OpenSQLite opens the database at path (":memory:" for an in-memory DB) and creates the schema.
func OpenSQLite(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	// SQLite allows only one writer; a single connection also keeps ":memory:" databases shared.
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return &SQLiteStore{db: db, now: func() time.Time { return time.Now().UTC() }}, nil
}

// Close closes the underlying database.
func (s *SQLiteStore) Close() error { return s.db.Close() }

const selectCols = `SELECT id, title, completed, created_at, updated_at FROM todos`

type scanner interface{ Scan(dest ...any) error }

func scanTodo(sc scanner) (Todo, error) {
	var t Todo
	var created, updated string
	if err := sc.Scan(&t.ID, &t.Title, &t.Completed, &created, &updated); err != nil {
		return Todo{}, err
	}
	var err error
	if t.CreatedAt, err = time.Parse(timeLayout, created); err != nil {
		return Todo{}, err
	}
	if t.UpdatedAt, err = time.Parse(timeLayout, updated); err != nil {
		return Todo{}, err
	}
	return t, nil
}

// List returns all todos, newest first.
func (s *SQLiteStore) List(ctx context.Context) ([]Todo, error) {
	rows, err := s.db.QueryContext(ctx, selectCols+` ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	todos := []Todo{}
	for rows.Next() {
		t, err := scanTodo(rows)
		if err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

// Get returns the todo with the given ID, or ErrNotFound.
func (s *SQLiteStore) Get(ctx context.Context, id int64) (Todo, error) {
	t, err := scanTodo(s.db.QueryRowContext(ctx, selectCols+` WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Todo{}, ErrNotFound
	}
	return t, err
}

// Create inserts a new todo after validating its title.
func (s *SQLiteStore) Create(ctx context.Context, title string) (Todo, error) {
	title, err := NormalizeTitle(title)
	if err != nil {
		return Todo{}, err
	}
	ts := s.now().Format(timeLayout)
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO todos (title, completed, created_at, updated_at) VALUES (?, 0, ?, ?)`,
		title, ts, ts)
	if err != nil {
		return Todo{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Todo{}, err
	}
	return s.Get(ctx, id)
}

// Update applies the non-nil fields of p to the todo with the given ID.
func (s *SQLiteStore) Update(ctx context.Context, id int64, p Patch) (Todo, error) {
	t, err := s.Get(ctx, id)
	if err != nil {
		return Todo{}, err
	}
	if p.Title != nil {
		if t.Title, err = NormalizeTitle(*p.Title); err != nil {
			return Todo{}, err
		}
	}
	if p.Completed != nil {
		t.Completed = *p.Completed
	}
	_, err = s.db.ExecContext(ctx,
		`UPDATE todos SET title = ?, completed = ?, updated_at = ? WHERE id = ?`,
		t.Title, t.Completed, s.now().Format(timeLayout), id)
	if err != nil {
		return Todo{}, err
	}
	return s.Get(ctx, id)
}

// Delete removes the todo with the given ID, or returns ErrNotFound.
func (s *SQLiteStore) Delete(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM todos WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
