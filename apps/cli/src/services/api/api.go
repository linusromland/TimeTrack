package apiService

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"TimeTrack-cli/src/database"
)

type APIService struct {
	db      *database.DBWrapper
	baseURL string
	client  *http.Client
}

func NewAPIService(db *database.DBWrapper) *APIService {
	if db == nil {
		panic("database wrapper cannot be nil")
	}

	rawURL := db.Get(database.ServerURLKey)
	if rawURL == "" {
		panic("server URL cannot be empty")
	}

	normalized := normalizeBaseURL(rawURL)

	return &APIService{
		db:      db,
		baseURL: normalized,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func normalizeBaseURL(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		panic(fmt.Sprintf("invalid server URL: %q", raw))
	}

	cleanPath := strings.TrimRight(parsed.Path, "/")
	parsed.Path = fmt.Sprintf("%s/api/v1", cleanPath)

	return parsed.String()
}

func (api *APIService) GetBaseURL() string {
	return api.baseURL
}

func (api *APIService) newAuthRequest(method, url string, body []byte) (*http.Request, error) {
	token := api.db.Get(database.AuthTokenKey)
	if token == "" {
		return nil, fmt.Errorf("no authentication token found")
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
