package usecase

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app/internal/domains/ihsg"
)

func TestGetMarketsKeepsSuccessfulCacheAfterProviderFailure(t *testing.T) {
	providerAvailable := true
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if !providerAvailable {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"chart":{"result":[{"timestamp":[1787616000],"indicators":{"quote":[{"open":[1],"high":[2],"low":[1],"close":[2],"volume":[10]}]}}],"error":null}}`)
	}))
	defer provider.Close()
	t.Setenv("IHSG_API_URL", provider.URL)

	originalConfigs := marketConfigs
	marketConfigs = append([]marketConfig(nil), marketConfigs...)
	for i := range marketConfigs {
		marketConfigs[i].StooqSymbol = ""
	}
	defer func() { marketConfigs = originalConfigs }()

	uc := NewIhsgUsecase().(*ihsgUsecase)
	markets, err := uc.GetMarkets(t.Context())
	if err != nil {
		t.Fatalf("first refresh failed: %v", err)
	}
	if len(markets[0].Data) != 1 {
		t.Fatal("expected the successful provider refresh to populate the cache")
	}
	if remaining := time.Until(uc.cacheUntil); remaining < successfulRefreshTTL-time.Second {
		t.Fatalf("expected successful refresh TTL, got %s", remaining)
	}

	providerAvailable = false
	uc.cacheUntil = time.Now().Add(-time.Second)
	markets, err = uc.GetMarkets(t.Context())
	if err != nil {
		t.Fatalf("failed refresh returned an error: %v", err)
	}
	if len(markets[0].Data) != 1 {
		t.Fatal("expected the last known good data after provider failure")
	}
	if remaining := time.Until(uc.cacheUntil); remaining < failedRefreshTTL-time.Second || remaining > failedRefreshTTL {
		t.Fatalf("expected failed refresh TTL, got %s", remaining)
	}
}

func TestPreserveCachedMarketsUsesLastKnownGoodData(t *testing.T) {
	cachedPoint := &ihsg.IHSGPoint{Date: "2026-08-24", Open: 1, High: 2, Low: 1, Close: 2, Volume: 10}
	cachedMarkets := []ihsg.MarketDataset{
		{
			Code:      "ID",
			Source:    "Yahoo Finance (IHSG)",
			FetchedAt: "2026-08-24T00:00:00Z",
			Data:      []*ihsg.IHSGPoint{cachedPoint},
		},
	}
	markets := []ihsg.MarketDataset{
		{
			Code:         "ID",
			FetchedAt:    "2026-08-25T00:00:00Z",
			Data:         []*ihsg.IHSGPoint{},
			ErrorMessage: "provider unavailable",
		},
	}

	hasFreshData := preserveCachedMarkets(markets, cachedMarkets)

	if hasFreshData {
		t.Fatal("expected a failed refresh to report no fresh data")
	}
	if len(markets[0].Data) != 1 || markets[0].Data[0] != cachedPoint {
		t.Fatal("expected the last known good market points to be preserved")
	}
	if markets[0].FetchedAt != cachedMarkets[0].FetchedAt || markets[0].Source != cachedMarkets[0].Source {
		t.Fatal("expected stale data to preserve its original source and fetch time")
	}
	if markets[0].ErrorMessage != "" {
		t.Fatal("expected the cached successful dataset to replace the failed refresh")
	}
}

func TestPreserveCachedMarketsKeepsFreshData(t *testing.T) {
	freshPoint := &ihsg.IHSGPoint{Date: "2026-08-25", Open: 2, High: 3, Low: 2, Close: 3, Volume: 20}
	markets := []ihsg.MarketDataset{
		{Code: "ID", FetchedAt: "2026-08-25T00:00:00Z", Data: []*ihsg.IHSGPoint{freshPoint}},
	}
	cachedMarkets := []ihsg.MarketDataset{
		{Code: "ID", FetchedAt: "2026-08-24T00:00:00Z", Data: []*ihsg.IHSGPoint{{Date: "2026-08-24"}}},
	}

	hasFreshData := preserveCachedMarkets(markets, cachedMarkets)

	if !hasFreshData {
		t.Fatal("expected a usable refresh to report fresh data")
	}
	if len(markets[0].Data) != 1 || markets[0].Data[0] != freshPoint {
		t.Fatal("expected fresh data to remain unchanged")
	}
}

func TestPreserveCachedMarketsDoesNotReuseAnotherMarket(t *testing.T) {
	markets := []ihsg.MarketDataset{{Code: "ID", Data: []*ihsg.IHSGPoint{}}}
	cachedMarkets := []ihsg.MarketDataset{
		{Code: "US", Data: []*ihsg.IHSGPoint{{Date: "2026-08-24"}}},
	}

	preserveCachedMarkets(markets, cachedMarkets)

	if len(markets[0].Data) != 0 {
		t.Fatal("expected cached data to be matched by market code")
	}
}
