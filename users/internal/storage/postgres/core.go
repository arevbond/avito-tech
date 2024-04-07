package postgres

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"strconv"
	"strings"
	"users/cmd/users/config"
)

const driverName = "pgx"

type Storage struct {
	db  *sqlx.DB
	log *slog.Logger
}

func New(log *slog.Logger, cfg config.StorageConfig) (*Storage, error) {
	connectionString := createDBConnectionString(cfg)
	parsedConnConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("can't parse config: %w", err)
	}
	db := sqlx.NewDb(stdlib.OpenDB(*parsedConnConfig), driverName)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't ping storage: %w", err)
	}

	return &Storage{
		db:  db,
		log: log,
	}, nil
}

func createDBConnectionString(cfg config.StorageConfig) string {
	connectionMap := map[string]string{
		"host":     cfg.Host,
		"port":     strconv.Itoa(cfg.Port),
		"database": cfg.Database,
		"user":     cfg.Username,
		"password": cfg.Password,
	}
	connectionSlice := make([]string, 0, len(connectionMap))
	for k, v := range connectionMap {
		connectionSlice = append(connectionSlice, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(connectionSlice, " ")
}
