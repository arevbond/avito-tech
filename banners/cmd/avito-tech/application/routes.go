package application

import (
	"banners/internal/publicapi"
	"github.com/go-chi/chi/v5"
)

func (a *App) newHTTPServer(env *env) *HTTPServerWrap {
	return NewHTTPServerWrap(
		a.Log,
		withAdminServer(a.Config.AdminServer),
		withPublicServer(a.Config.PublicServer, a.publicMux(env)))
}

func (a *App) publicMux(env *env) *chi.Mux {
	mux := chi.NewMux()

	handler := publicapi.Handler{
		Log:           a.Log,
		BannerService: env.bannerService,
	}

	mux.Get("/user_banner", handler.UserBanner)
	//mux.Get("/banner", handler.)
	mux.Post("/banner", handler.CreateBanner)
	mux.Patch("/banner/{id}", handler.UpdateBanner)
	mux.Delete("/banner/{id}", handler.DeleteBanner)

	//mux.Route("/banner", func(r chi.Router) {
	//r.Use(bannerCtx)
	//r.Patch("/{id}", handler.UpdateBanner)
	//	r.Delete("/{id}", .handler)
	//})

	return mux
}
