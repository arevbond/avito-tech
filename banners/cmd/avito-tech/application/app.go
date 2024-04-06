package application

import (
	"banners/cmd/avito-tech/config"
	"banners/internal/service/banner"
	"banners/internal/storage/postgres"
	"context"
	"fmt"
	"log/slog"
	"syscall"
)

type App struct {
	Config *config.Config
	Log    *slog.Logger
	Closer *Closer
}

func New(cfg *config.Config, log *slog.Logger) *App {
	return &App{
		Config: cfg,
		Log:    log,
		Closer: NewCloser(log, cfg.Application.GracefulShutdownTimeout, syscall.SIGINT, syscall.SIGTERM),
	}
}

func (a *App) Run() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	a.Closer.Add(func() error {
		cancelFunc()
		return nil
	})

	envStruct, err := a.constructEnv(ctx)
	if err != nil {
		return fmt.Errorf("can't constuct enviroment: %w", err)
	}

	httpServer := a.newHTTPServer(envStruct)
	a.Closer.Add(httpServer.GracefulStop()...)

	a.Closer.Run(httpServer.Run()...)
	a.Closer.Wait()
	return nil
}

type env struct {
	bannerService banner.Service
}

func (a *App) constructEnv(ctx context.Context) (*env, error) {
	db, err := postgres.New(a.Log, a.Config.Storage)
	if err != nil {
		return nil, fmt.Errorf("can't init storage: %w", err)
	}

	a.Closer.Add(db.Close)

	bannerService := &banner.ServiceImpl{
		Storage: db,
	}

	return &env{
		bannerService: bannerService,
	}, nil
}
