package database

import "github.com/jackc/pgx/v4"

type Row interface {
	Scan(dest ...interface{}) error
}

type databaseRow struct {
	row pgx.Row
}

func newDatabaseRow(row pgx.Row) Row {
	return &databaseRow{
		row: row,
	}
}

func (p *databaseRow) Scan(dest ...interface{}) error {
	return p.row.Scan(dest...)
}
