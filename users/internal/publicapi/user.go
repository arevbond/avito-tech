package publicapi

import (
	"errors"
	"fmt"
	"net/http"
	"users/internal/service"
	"users/internal/utils"
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
			IsAdmin:  request.IsAdmin,
		})
		if err != nil {
			return nil, utils.WrapError(err, "user already exist", 400)
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

type tokenRequest struct {
	Token string `json:"token"`
}

func (h *Handler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	handlerRequest := func() (*verifyTokenResponse, error) {
		ctx := r.Context()

		request, err := parseJSONRequest[tokenRequest](r)
		if err != nil {
			return nil, fmt.Errorf("can't parse json: %w", err)
		}
		result, err := h.userService.VerifyToken(ctx, &service.TokenParams{Token: request.Token})
		if err != nil {
			return nil, fmt.Errorf("verify token service: %w", err)
		}
		return &verifyTokenResponse{Valid: result}, nil
	}
	response, err := handlerRequest()
	if err != nil {
		h.writeError(w, fmt.Errorf("verify token: %w", err))
		return
	}
	h.writeResponse(w, response)
}

func (h *Handler) IsAdmin(w http.ResponseWriter, r *http.Request) {
	handlerRequest := func() (*isAdminResponse, error) {
		ctx := r.Context()

		request, err := parseJSONRequest[tokenRequest](r)
		if err != nil {
			return nil, fmt.Errorf("can't parse json: %w", err)
		}

		isValid, err := h.userService.VerifyToken(ctx, &service.TokenParams{Token: request.Token})
		if err != nil {
			return nil, fmt.Errorf("verify token service: %w", err)
		}
		if !isValid {
			return nil, utils.WrapError(errors.New("invalid token"), "invalid token", 404)
		}

		result, err := h.userService.IsAdmin(ctx, &service.TokenParams{Token: request.Token})
		if err != nil {
			return nil, fmt.Errorf("verify token service: %w", err)
		}
		return &isAdminResponse{IsAdmin: result}, nil
	}
	response, err := handlerRequest()
	if err != nil {
		h.writeError(w, fmt.Errorf("is admin token: %w", err))
		return
	}
	h.writeResponse(w, response)
}
