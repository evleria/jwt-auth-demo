package database

import (
	"context"
	//"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database interface {
	Exec(sql string, args ...interface{}) (int64, error)
}

type database struct {
	ctx  context.Context
	conn *pgxpool.Pool
}

func New(ctx context.Context, connString string) (Database, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &database{
		ctx:  ctx,
		conn: conn,
	}, nil
}

func (p *database) Exec(sql string, args ...interface{}) (int64, error) {
	result, err := p.conn.Exec(p.ctx, sql, args...)
	return result.RowsAffected(), err
}
