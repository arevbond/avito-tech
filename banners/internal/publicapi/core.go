package publicapi

import (
	"banners/internal/models"
	"banners/internal/service"
	"banners/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Handler struct {
	Log           *slog.Logger
	BannerService service.Service
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	errorResult, ok := utils.FromError(err)
	if !ok {
		h.Log.Error("can't write error in response")
		return
	}
	w.WriteHeader(errorResult.StatusCode)
	err = json.NewEncoder(w).Encode(
		map[string]any{
			"error": errorResult.Msg,
		})
	if err != nil {
		http.Error(w, utils.InternalErrorMessage, http.StatusInternalServerError)
	}
}

func (h *Handler) writeResponse(w http.ResponseWriter, response any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if statusCode == http.StatusNoContent {
		return
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, utils.InternalErrorMessage, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) parseJSONBody(r *http.Request) (*models.CreateBanner, error) {
	var banner *models.CreateBanner
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read request body: %w", err)
	}
	err = json.Unmarshal(data, &banner)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal json: %w", err)
	}
	return banner, nil
}
