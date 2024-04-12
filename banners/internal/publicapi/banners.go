package publicapi

import (
	"banners/internal/models"
	"banners/internal/service"
	"banners/internal/utils"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type userBannerResponse struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

func (h *Handler) UserBanner(w http.ResponseWriter, r *http.Request) {
	handleRequst := func() (*userBannerResponse, *utils.ErrorResult) {
		ctx := r.Context()

		tagIDstr := r.URL.Query().Get("tag_id")
		featureIDstr := r.URL.Query().Get("feature_id")

		if tagIDstr == "" || featureIDstr == "" {
			return nil, utils.ErrNotRequiredParam
		}

		tagID, err := strconv.Atoi(tagIDstr)
		if err != nil {
			h.Log.Error("can't convert string", "error", err)
			return nil, utils.ErrInvalidTypeParam
		}

		featureID, err := strconv.Atoi(featureIDstr)
		if err != nil {
			h.Log.Error("can't convert string", "error", err)
			return nil, utils.ErrInvalidTypeParam
		}

		var defaultUseLastRevision bool
		useLastRevisionString := r.URL.Query().Get("use_last_revision")

		useLastRevision, err := strconv.ParseBool(useLastRevisionString)
		if err != nil {
			h.Log.Error("can't convert string", "error", err)
			useLastRevision = defaultUseLastRevision
		}

		token := r.Header.Get("token")
		if token == "" {
			return nil, utils.WrapServiceError(service.ErrUnauthorized)
		}

		userBanner, err := h.BannerService.UserBanner(ctx, token, &service.UserBannerParams{
			TagID:           tagID,
			FeatureID:       featureID,
			UseLastRevision: useLastRevision,
		})

		if err != nil {
			if errors.Is(service.ErrNotFound, err) {
				return nil, utils.WrapServiceError(err)
			}
			return nil, utils.WrapInternalError(fmt.Errorf("service get banner: %w", err))
		}

		return &userBannerResponse{
			Title: userBanner.Title,
			Text:  userBanner.Text,
			Url:   userBanner.Url,
		}, nil
	}

	response, err := handleRequst()
	if err != nil {
		h.writeError(w, fmt.Errorf("user banner: %w", err))
		return
	}
	h.writeResponse(w, response, http.StatusOK)
}

func (h *Handler) Banners(w http.ResponseWriter, r *http.Request) {
	handleRequst := func() ([]*models.Banner, *utils.ErrorResult) {
		ctx := r.Context()

		defaultTagID := 0
		defaultFeatureID := 0
		defaultOffset := 0
		defaultLimit := 100

		tagIDstr := r.URL.Query().Get("tag_id")
		featureIDstr := r.URL.Query().Get("feature_id")
		offsetStr := r.URL.Query().Get("offset")
		limitStr := r.URL.Query().Get("limit")

		tagID, err := strconv.Atoi(tagIDstr)
		if err != nil {
			tagID = defaultTagID
		}

		featureID, err := strconv.Atoi(featureIDstr)
		if err != nil {
			featureID = defaultFeatureID
		}

		// Validate and parse offset
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = defaultOffset
		}

		// Validate and parse limit
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = defaultLimit
		}

		token := r.Header.Get("token")
		if token == "" {
			return nil, utils.WrapServiceError(service.ErrUnauthorized)
		}

		banners, err := h.BannerService.Banners(ctx, token, &models.BannersParams{
			TagID:     tagID,
			FeatureID: featureID,
			Offset:    offset,
			Limit:     limit,
		})
		if err != nil {
			h.Log.Error("banners error", "error", err)
			return nil, utils.WrapServiceError(err)
		}
		return banners, nil
	}

	response, err := handleRequst()
	if err != nil {
		h.writeError(w, fmt.Errorf("banners: %w", err))
		return
	}
	h.writeResponse(w, response, http.StatusOK)
}

type createBannerResponse struct {
	BannerID int `json:"banner_id"`
}

func (h *Handler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() (*createBannerResponse, *utils.ErrorResult) {
		ctx := r.Context()

		token := r.Header.Get("token")
		if token == "" {
			return nil, utils.WrapServiceError(service.ErrUnauthorized)
		}

		banner, err := h.parseJSONBody(r)
		if err != nil {
			return nil, utils.WrapInternalError(err)
		}

		id, err := h.BannerService.CreateBanner(ctx, token, banner)
		if err != nil {
			return nil, utils.WrapServiceError(err)
		}
		return &createBannerResponse{BannerID: id}, nil
	}

	response, err := handleRequest()
	if err != nil {
		h.writeError(w, fmt.Errorf("create banner: %w", err))
		return
	}
	h.writeResponse(w, response, http.StatusCreated)
}

func (h *Handler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() error {
		ctx := r.Context()

		token := r.Header.Get("token")
		if token == "" {
			return utils.WrapServiceError(service.ErrUnauthorized)
		}

		var bannerIDstr string
		if bannerIDstr = chi.URLParam(r, "id"); bannerIDstr == "" {
			return utils.ErrNotRequiredParam
		}
		bannerID, err := strconv.Atoi(bannerIDstr)
		if err != nil {
			return utils.ErrInvalidTypeParam
		}

		banner, err := h.parseJSONBody(r)
		if err != nil {
			return utils.WrapInternalError(err)
		}

		err = h.BannerService.UpdateBanner(ctx, token, bannerID, banner)
		if err != nil {
			return utils.WrapServiceError(err)
		}
		return nil
	}

	err := handleRequest()
	if err != nil {
		h.writeError(w, fmt.Errorf("update banner: %w", err))
		return
	}
	h.writeResponse(w, map[string]string{"status": "OK"}, http.StatusOK)
}

func (h *Handler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() error {
		ctx := r.Context()

		token := r.Header.Get("token")
		if token == "" {
			return utils.WrapServiceError(service.ErrUnauthorized)
		}

		var bannerIDstr string
		if bannerIDstr = chi.URLParam(r, "id"); bannerIDstr == "" {
			return utils.ErrNotRequiredParam
		}
		bannerID, err := strconv.Atoi(bannerIDstr)
		if err != nil {
			return utils.ErrInvalidTypeParam
		}

		err = h.BannerService.DeleteBanner(ctx, token, bannerID)
		if err != nil {
			return fmt.Errorf("service error: %w", err)
		}
		return nil
	}

	err := handleRequest()
	if err != nil {
		h.writeError(w, fmt.Errorf("update banner: %w", err))
		return
	}
	h.writeResponse(w, map[string]string{"status": "OK"}, http.StatusNoContent)
}
