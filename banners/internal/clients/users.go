package clients

import (
	"banners/cmd/avito-tech/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	verifyTokenEndpoint  = "verify-token"
	isAdminTokenEndpoint = "is-admin"
	schema               = "http"
)

type Users struct {
	address string
	port    int
	client  http.Client
}

func New(cfg config.UsersServiceConfig) *Users {
	return &Users{
		address: cfg.Address,
		port:    cfg.Port,
		client:  http.Client{},
	}
}

func (u *Users) doRequest(data []byte, method string, url string) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("can't create new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't convert response body to bytes: %w", err)
	}

	return result, nil
}

func (u *Users) formatURL(endpoint string) string {
	return fmt.Sprintf("%s://%s:%d/%s", schema, u.address, u.port, endpoint)
}

type tokenParams struct {
	Token string `json:"token"`
}

type verifyTokenResponse struct {
	Valid bool `json:"valid"`
}

type isAdminTokenResponse struct {
	IsAdmin bool `json:"is_admin"`
}

func (u *Users) getParams(token string) ([]byte, error) {
	params := tokenParams{Token: token}
	data, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("can't marshal token: %w", err)
	}
	return data, nil
}

func (u *Users) VerifyToken(token string) (bool, error) {
	data, err := u.getParams(token)
	if err != nil {
		return false, fmt.Errorf("can't marshal token: %w", err)
	}
	resp, err := u.doRequest(data, http.MethodPost, u.formatURL(verifyTokenEndpoint))
	if err != nil {
		return false, fmt.Errorf("can't get response: %w", err)
	}
	var result verifyTokenResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return false, fmt.Errorf("can't unmarshal json: %w", err)
	}
	return result.Valid, nil
}

func (u *Users) IsAdmin(token string) (bool, error) {
	data, err := u.getParams(token)
	if err != nil {
		return false, fmt.Errorf("can't marshal token: %w", err)
	}
	resp, err := u.doRequest(data, http.MethodPost, u.formatURL(isAdminTokenEndpoint))
	if err != nil {
		return false, fmt.Errorf("can't get response: %w", err)
	}
	var result isAdminTokenResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return false, fmt.Errorf("can't unmarshal json: %w", err)
	}
	return result.IsAdmin, nil
}
