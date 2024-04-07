package publicapi

import (
	"fmt"
	"net/http"
	"users/internal/service"
)

type jwtTokenResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

type verifyTokenResponse struct {
	Valid bool `json:"valid"`
}

type isAdminResponse struct {
	IsAdmin bool `json:"is_admin"`
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() (*jwtTokenResponse, error) {
		ctx := r.Context()

		request, err := parseJSONRequest[registerRequest](r)
		if err != nil {
			return nil, fmt.Errorf("can't parse json: %w", err)
		}

		tokenModel, err := h.userService.Register(ctx, &service.RegisterParams{
			Username: request.Username,
			Password: request.Password,
		})
		if err != nil {
			return nil, fmt.Errorf("register service: %w", err)
		}

		return &jwtTokenResponse{
			Token:  tokenModel.Value,
			UserID: tokenModel.UserID.String(),
		}, nil
	}

	response, err := handleRequest()
	if err != nil {
		h.writeError(w, fmt.Errorf("register: %w", err))
		return
	}
	h.writeResponse(w, response)
}
