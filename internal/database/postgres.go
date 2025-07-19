package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"rinha2025/internal/config"

	_ "github.com/lib/pq"
)

type (
	Database struct {
		Connection *sql.DB
	}

	InitOptions struct {
		Host         string
		Port         int
		User         string
		Password     string
		DatabaseName string
	}
)

const ConnectionString = "postgres://%s:%s@%s:%d/%s?sslmode=disable"

func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	connStr := fmt.Sprintf(ConnectionString,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
	)

	conn, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(cfg.MaxIdleConnections)
	conn.SetMaxOpenConns(cfg.MaxOpenConnections)

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &Database{
		Connection: conn,
	}, nil
}

func (d *Database) Close() {
	if err := d.Connection.Close(); err != nil {
		slog.Error("error closing database connection", slog.String("error", err.Error()))
	}
}
