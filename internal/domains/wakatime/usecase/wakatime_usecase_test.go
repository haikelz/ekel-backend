package usecase

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatsSendsOAuthTokenAndRange(t *testing.T) {
	wakatime := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer waka_tok_test" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"error":"Unauthorized"}`)
			return
		}
		if r.URL.Path != "/all_time" {
			t.Errorf("expected /all_time, got %s", r.URL.Path)
		}
		fmt.Fprint(w, `{"data":{"human_readable_total":"1 hr","range":"all_time"}}`)
	}))
	defer wakatime.Close()
	t.Setenv("WAKATIME_API_URL", wakatime.URL)
	uc := NewWakatimeUsecase()

	t.Setenv("WAKATIME_ACCESS_TOKEN", "waka_tok_test")
	stats, err := uc.GetStats(t.Context(), "all_time")
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	if stats.HumanReadableTotal != "1 hr" {
		t.Fatalf("unexpected stats: %+v", stats)
	}

	t.Setenv("WAKATIME_ACCESS_TOKEN", "waka_legacy_api_key")
	if _, err := uc.GetStats(t.Context(), "all_time"); err == nil {
		t.Fatal("expected an error when WakaTime rejects the token")
	}

	if _, err := uc.GetStats(t.Context(), "last_day"); err == nil {
		t.Fatal("expected an error for an undocumented range")
	}
}
