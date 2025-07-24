package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"rinha2025/internal/config"

	_ "github.com/lib/pq"
)

type Database struct {
	Connection *sql.DB
}

const ConnectionString = "postgres://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s"

func NewDatabase(cfg config.DatabaseConfig, appName string) (*Database, error) {
	connStr := fmt.Sprintf(ConnectionString,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
		appName,
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
