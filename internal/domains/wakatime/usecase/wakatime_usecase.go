package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"app/internal/domains/wakatime"
)

// Ranges documented at https://wakatime.com/developers#stats. Without a range,
// WakaTime falls back to the public profile range, so always send one.
var validRanges = map[string]bool{
	"last_7_days":   true,
	"last_30_days":  true,
	"last_6_months": true,
	"last_year":     true,
	"all_time":      true,
}

type wakatimeUsecase struct {
	client *http.Client
}

func NewWakatimeUsecase() wakatime.WakatimeUsecase {
	return &wakatimeUsecase{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (uc *wakatimeUsecase) GetStats(ctx context.Context, rng string) (*wakatime.StatsResponse, error) {
	if !validRanges[rng] {
		return nil, fmt.Errorf("unsupported range %q", rng)
	}

	baseURL := os.Getenv("WAKATIME_API_URL")
	if baseURL == "" {
		baseURL = "https://wakatime.com/api/v1/users/current/stats"
	}

	// Since 2026-11-01 WakaTime API keys only work for plugins; reads need a personal OAuth token.
	token := os.Getenv("WAKATIME_ACCESS_TOKEN")
	if token == "" {
		return nil, errors.New("WAKATIME_ACCESS_TOKEN not set")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/"+rng, nil)
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

	// 202 means WakaTime is still refreshing the stats; keep whatever data it returned.
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("wakatime returned %d: %s", res.StatusCode, body)
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
