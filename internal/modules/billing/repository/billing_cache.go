package billingrepository

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
)

type BillingStoreCache struct {
	rdb *redis.Client
}

func NewBillingCacheStore(rdb *redis.Client) *BillingStoreCache {
	return &BillingStoreCache{rdb: rdb}
}

func (r *BillingStoreCache) SetRates(ctx context.Context, rates map[string]float64) error {
	ratesMap := make(map[string]interface{}, len(rates))
	for k, v := range rates {
		ratesMap[k] = v
	}

	pipe := r.rdb.Pipeline()
	pipe.HSet(ctx, billingdomain.LatestRatesKey, ratesMap)
	pipe.Expire(ctx, billingdomain.LatestRatesKey, billingdomain.LatestRatesExpiration)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *BillingStoreCache) GetRates(ctx context.Context) (map[string]float64, error) {
	stringRates, err := r.rdb.HGetAll(ctx, billingdomain.LatestRatesKey).Result()
	if err != nil {
		return nil, err
	}

	if len(stringRates) == 0 {
		return nil, redis.Nil
	}

	rates := make(map[string]float64, len(stringRates))
	for currency, rateStr := range stringRates {
		rateFloat, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			log.Printf("Error al parsear tasa para %s: %v", currency, err)
			continue
		}
		rates[currency] = rateFloat
	}

	return rates, nil
}

func (r *BillingStoreCache) DeleteRates(ctx context.Context) error {
	return r.rdb.Del(ctx, billingdomain.LatestRatesKey).Err()
}

func (r *BillingStoreCache) SetSpecificTimeRateForCurrency(ctx context.Context, date string, currency string, rate float64) error {
	key := fmt.Sprintf("%s:%s:%s", billingdomain.RatePrefix, date, currency)

	return r.rdb.Set(ctx, key, rate, billingdomain.RateExpiration).Err()
}

func (r *BillingStoreCache) GetSpecificTimeRateForCurrency(ctx context.Context, date string, currency string) (float64, error) {
	key := fmt.Sprintf("%s:%s:%s", billingdomain.RatePrefix, date, currency)

	rateStr, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(rateStr, 64)
}
