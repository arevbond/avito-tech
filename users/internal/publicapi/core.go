package publicapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log/slog"
	"net/http"
	"users/internal/service"
	"users/internal/storage"
	"users/internal/utils"
)

type Handler struct {
	log         *slog.Logger
	userService service.Service
}

func Routes(log *slog.Logger, db storage.Storage) chi.Router {
	r := chi.NewRouter()
	handler := &Handler{log, service.New(db, log)}
	r.Post("/register", handler.Register)
	r.Post("/verify-token", handler.VerifyToken)
	r.Post("/is-admin", handler.IsAdmin)
	return r
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	h.log.Error("http response error", "error", err)

	w.Header().Set("Content-Type", "application/json")
	errorResult, ok := utils.FromError(err)
	if !ok {
		h.log.Error("can't write log message")
		return
	}
	w.WriteHeader(errorResult.StatusCode)
	err = json.NewEncoder(w).Encode(
		map[string]any{
			"message": errorResult.Msg,
			"code":    errorResult.StatusCode,
		})

	if err != nil {
		http.Error(w, utils.InternalErrorMessage, http.StatusInternalServerError)
	}
}

func (h *Handler) writeResponse(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, utils.InternalErrorMessage, http.StatusInternalServerError)
	}
}

func parseJSONRequest[T registerRequest | tokenRequest](r *http.Request) (*T, error) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("can't read body from request: %w", err)
		return nil, utils.WrapInternalError(err)
	}

	var request T
	err = json.Unmarshal(body, &request)
	if err != nil {
		err = fmt.Errorf("can't unmarshall json: %w", err)
		return nil, utils.WrapInternalError(err)
	}
	return &request, nil
}
