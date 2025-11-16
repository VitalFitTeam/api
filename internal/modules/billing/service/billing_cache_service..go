package billingservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

type OpenExchangeResponse struct {
	Disclaimer string             `json:"disclaimer"`
	License    string             `json:"license"`
	Timestamp  int64              `json:"timestamp"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}

func (bs *BillingService) GetLatestRates(ctx context.Context) (map[string]float64, error) {
	cachedRates, err := bs.cache.Billing.GetRates(ctx)
	if err == nil {
		log.Println("CACHE HIT: Getting 'latest_rates' from Redis.")
		return cachedRates, nil
	}

	if !errors.Is(err, redis.Nil) {
		log.Printf("Redis error getting 'latest_rates': %v. Falling back to API.", err)
	}

	log.Println("CACHE MISS: Getting 'latest_rates' from the API.")
	apiRates, err := bs.fetchLatestRatesFromAPI(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from API: %w", err)
	}

	go func() {
		log.Println("CACHE SET: Saving 'latest_rates' to Redis in the background.")
		err := bs.cache.Billing.SetRates(context.Background(), apiRates)
		if err != nil {
			log.Printf("Error caching 'latest_rates': %v", err)
		}
	}()

	return apiRates, nil
}

func (bs *BillingService) GetHistoricalRateForCurrency(ctx context.Context, date, currency string) (float64, error) {

	keyName := fmt.Sprintf("%s:%s", date, currency)
	cachedRate, err := bs.cache.Billing.GetSpecificTimeRateForCurrency(ctx, date, currency)
	if err == nil {
		log.Printf("CACHE HIT: Getting 'rate:%s' from Redis.", keyName)
		return cachedRate, nil
	}

	if !errors.Is(err, redis.Nil) {
		log.Printf("Redis error getting 'rate:%s': %v. Falling back to API.", keyName, err)
	}

	log.Printf("CACHE MISS: Getting 'rate:%s' from the API.", keyName)
	apiRate, err := bs.fetchHistoricalRateFromAPI(ctx, date, currency)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch from API: %w", err)
	}

	go func() {
		log.Printf("CACHE SET: Saving 'rate:%s' to Redis (expires in 3 min).", keyName)
		err := bs.cache.Billing.SetSpecificTimeRateForCurrency(context.Background(), date, currency, apiRate)
		if err != nil {
			log.Printf("Error caching 'rate:%s': %v", keyName, err)
		}
	}()

	return apiRate, nil
}

func (bs *BillingService) fetchLatestRatesFromAPI(ctx context.Context) (map[string]float64, error) {
	appID := bs.cfg.OpenExchange.AppID
	if appID == "" {
		return nil, errors.New("OpenExchange AppID is not configured")
	}

	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", appID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := bs.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned a non-OK status: %s", resp.Status)
	}

	var apiResponse OpenExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return apiResponse.Rates, nil
}

func (bs *BillingService) fetchHistoricalRateFromAPI(ctx context.Context, date, currency string) (float64, error) {
	appID := bs.cfg.OpenExchange.AppID
	if appID == "" {
		return 0, errors.New("OpenExchange AppID is not configured")
	}

	url := fmt.Sprintf("https://openexchangerates.org/api/historical/%s.json?app_id=%s&symbols=%s", date, appID, currency)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := bs.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned a non-OK status: %s", resp.Status)
	}

	var apiResponse OpenExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return 0, err
	}

	rate, ok := apiResponse.Rates[currency]
	if !ok {
		return 0, fmt.Errorf("currency %s not found in API response for date %s", currency, date)
	}

	return rate, nil
}
