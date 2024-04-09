package publicapi

import (
	"banners/internal/service/banner"
	"log/slog"
)

type Handler struct {
	Log           *slog.Logger
	BannerService banner.Service
}
