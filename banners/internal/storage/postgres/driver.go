package postgres

import (
	"banners/cmd/avito-tech/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"golang.yandex/hasql"
	"golang.yandex/hasql/checkers"
	"log/slog"
	"strconv"
	"strings"
)

const driverName = "pgx"

type PGStorage struct {
	log     *slog.Logger
	cluster *hasql.Cluster
}

func newPGStorage(log *slog.Logger, cfg config.StorageConfig) (*PGStorage, error) {
	cluster, err := newPGCluster(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't init posgres cluster: %w", err)
	}

	return &PGStorage{
		log:     log,
		cluster: cluster,
	}, nil
}

func (s *PGStorage) Close() error {
	return s.cluster.Close()
}

func (s *PGStorage) Master() *sqlx.DB {
	db := s.cluster.Primary().DB()
	return sqlx.NewDb(db, driverName)
}

func (s *PGStorage) Slave() *sqlx.DB {
	db := s.cluster.StandbyPreferred().DB()
	return sqlx.NewDb(db, driverName)
}

func newPGCluster(cfg config.StorageConfig) (*hasql.Cluster, error) {
	nodes := make([]hasql.Node, 0, len(cfg.Hosts))
	for _, host := range cfg.Hosts {
		connString := createDBConnectionString(host, cfg)
		parsedConnConfig, err := pgx.ParseConfig(connString)
		if err != nil {
			return nil, fmt.Errorf("can't parse connection config: %w", err)
		}
		db := sqlx.NewDb(stdlib.OpenDB(*parsedConnConfig), driverName)
		nodes = append(nodes, hasql.NewNode(host, db.DB))
	}

	cluster, err := hasql.NewCluster(nodes, checkers.PostgreSQL)
	if err != nil {
		return nil, fmt.Errorf("can't init new cluster: %w", err)
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), cfg.InitializationTimeout)
	defer cancelFunc()
	_, err = cluster.WaitForPrimary(ctx)
	if err != nil {
		if closeError := cluster.Close(); closeError != nil {
			return nil, fmt.Errorf("cluster close error: %w", closeError)
		}
		return nil, fmt.Errorf("wait for primary timeout exceed: %w", err)
	}
	return cluster, nil
}

func createDBConnectionString(host string, cfg config.StorageConfig) string {
	connectionMap := map[string]string{
		"host":     host,
		"port":     strconv.Itoa(cfg.Port),
		"database": cfg.Database,
		"user":     cfg.Username,
		"password": cfg.Password,
	}
	if cfg.SSLMode != "" {
		connectionMap["sslmode"] = cfg.SSLMode
	}

	connectionSlice := make([]string, 0, len(connectionMap))
	for k, v := range connectionMap {
		connectionSlice = append(connectionSlice, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(connectionSlice, " ")
}
