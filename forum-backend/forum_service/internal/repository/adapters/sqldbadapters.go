package adapters

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type DbAdapter struct {
	*sql.DB
}

func (d *DbAdapter) Exec(query string, args ...any) (sql.Result, error) {
	return d.DB.Exec(query, args...)
}

func (d *DbAdapter) Get(dest interface{}, query string, args ...interface{}) error {
	return d.DB.QueryRow(query, args...).Scan(dest)
}

func (d *DbAdapter) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.DB.ExecContext(ctx, query, args...)
}

func (d *DbAdapter) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, err := d.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	return sqlx.StructScan(rows, dest)
}

func (d *DbAdapter) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.DB.QueryRowContext(ctx, query, args...)
}

func (d *DbAdapter) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.DB.QueryContext(ctx, query, args...)
}

func (d *DbAdapter) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.DB.QueryRowContext(ctx, query, args...).Scan(dest)
}
