package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"app/internal/domains/wakatime"
)

type wakatimeUsecase struct {
	client *http.Client
}

func NewWakatimeUsecase() wakatime.WakatimeUsecase {
	return &wakatimeUsecase{
		client: &http.Client{},
	}
}

func (uc *wakatimeUsecase) GetStats(ctx context.Context, rng string) (*wakatime.StatsResponse, error) {
	baseURL := os.Getenv("WAKATIME_API_URL")
	if baseURL == "" {
		baseURL = "https://wakatime.com/api/v1/users/current/stats"
	}

	apiURL := baseURL
	switch rng {
	case "last_7_days":
		apiURL = baseURL + "/last_7_days"
	case "last_day":
		apiURL = baseURL + "/last_day"
	}

	// Since 2026-11-01 WakaTime API keys only work for plugins; reads need a personal OAuth token.
	token := os.Getenv("WAKATIME_ACCESS_TOKEN")
	if token == "" {
		return nil, errors.New("WAKATIME_ACCESS_TOKEN not set")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := uc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wakatime stats: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var raw struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wakatime response: %w", err)
	}

	var stats wakatime.StatsResponse
	if err := json.Unmarshal(raw.Data, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wakatime stats: %w", err)
	}

	return &stats, nil
}
