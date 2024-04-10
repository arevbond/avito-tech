package publicapi

import (
	"banners/internal/service"
	"banners/internal/utils"
	"errors"
	"fmt"
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
			return nil, utils.WrapError(errors.New("required param is absent"),
				"required param is absent", http.StatusBadRequest)
		}

		tagID, err := strconv.Atoi(tagIDstr)
		if err != nil {
			return nil, utils.WrapError(errors.New("invalid type parametr"),
				"invalid type parametr", http.StatusBadRequest)
		}

		featureID, err := strconv.Atoi(featureIDstr)
		if err != nil {
			return nil, utils.WrapError(errors.New("invalid type parametr"),
				"invalid type parametr", http.StatusBadRequest)
		}

		useLastRevisionString := r.URL.Query().Get("user_last_revision")

		userLastRevision, err := strconv.ParseBool(useLastRevisionString)
		if err != nil {
			return nil, utils.WrapError(errors.New("invalid type parametr"),
				"invalid type parametr", http.StatusBadRequest)
		}

		token := r.Header.Get("token")
		if token == "" {
			return nil, utils.WrapServiceError(service.ErrUnauthorized)
		}

		userBanner, err := h.BannerService.UserBanner(ctx, &service.UserBannerParams{
			TagID:           tagID,
			FeatureID:       featureID,
			UseLastRevision: userLastRevision,
			Token:           token,
		})
		if err != nil {
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
	h.writeResponse(w, response)
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
	h.writeResponse(w, response)
}
