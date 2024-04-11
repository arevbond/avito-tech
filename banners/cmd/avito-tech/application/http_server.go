package application

import (
	"banners/cmd/avito-tech/config"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

type serverOption struct {
	Port int
	Mux  *chi.Mux
}

type httpServerOption struct {
	adminServerOption   *serverOption
	publicServerOptions []*serverOption
}

type HTTPServerOption func(option *httpServerOption)

func withAdminServer(cfg config.ServerConfig) HTTPServerOption {
	return func(opts *httpServerOption) {
		opts.adminServerOption = &serverOption{Port: cfg.Port}
	}
}

func withPublicServer(cfg config.ServerConfig, mux *chi.Mux) HTTPServerOption {
	return func(opts *httpServerOption) {
		opts.publicServerOptions = append(opts.publicServerOptions,
			&serverOption{Port: cfg.Port, Mux: mux})
	}
}

type HTTPServerWrap struct {
	log     *slog.Logger
	servers []*http.Server
}

func NewHTTPServerWrap(log *slog.Logger, opts ...HTTPServerOption) *HTTPServerWrap {
	options := &httpServerOption{
		adminServerOption:   nil,
		publicServerOptions: nil,
	}

	for _, o := range opts {
		o(options)
	}

	var servers []*http.Server

	if options.adminServerOption != nil {
		servers = append(servers, NewNetHTTPServer(log, options.adminServerOption.Port, nil))
	}

	for _, option := range options.publicServerOptions {
		servers = append(servers, NewNetHTTPServer(log, option.Port, option.Mux))
	}
	return &HTTPServerWrap{
		log:     log,
		servers: servers,
	}
}

func (h *HTTPServerWrap) Run() []func() error {
	runFunc := func(server *http.Server) error {
		h.log.Info("run http server", slog.String("address", server.Addr))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server can't listen and server: %w", err)
		}
		return nil
	}

	response := make([]func() error, 0, len(h.servers))
	for _, server := range h.servers {
		response = append(response, func() error { return runFunc(server) })
	}
	return response
}

func (h *HTTPServerWrap) GracefulStop() []func() error {
	gracefulFunc := func(server *http.Server) error {
		err := server.Shutdown(context.Background())
		if err != nil {
			return fmt.Errorf("server shutdown: %w", err)
		}
		return nil
	}

	response := make([]func() error, 0, len(h.servers))
	for _, server := range h.servers {
		response = append(response, func() error { return gracefulFunc(server) })
	}
	return response
}

func NewNetHTTPServer(log *slog.Logger, port int, incomeMux *chi.Mux) *http.Server {
	mux := chi.NewMux()
	mux.Use(loggerMiddleware(log))
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.URLFormat)
	if incomeMux != nil {
		mux.Mount("/", incomeMux)
	}
	mux.Get("/ping", pingHandler())

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}

func loggerMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("http request", slog.String("host", r.Host),
				slog.String("uri", r.RequestURI))
			next.ServeHTTP(w, r)
		})
		return fn
	}
}

func pingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK\n"))
	}
}
