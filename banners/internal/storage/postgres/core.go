package postgres

import (
	"banners/cmd/avito-tech/config"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type Storage struct {
	storage *PGStorage
}

func New(log *slog.Logger, cfg config.StorageConfig) (*Storage, error) {
	postgresStorage, err := newPGStorage(log, cfg)
	if err != nil {
		return nil, fmt.Errorf("can't init new pg storage: %w", err)
	}
	return &Storage{
		storage: postgresStorage,
	}, nil
}

func (s *Storage) Close() error {
	return s.storage.Close()
}

func (s *Storage) Master() sqlx.ExtContext {
	return s.storage.Master()
}

func (s *Storage) Slave() sqlx.ExtContext {
	return s.storage.Slave()
}
