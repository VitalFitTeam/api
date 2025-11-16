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

// SetRates implementa el almacenamiento de tasas usando un Hash de Redis con una expiración de 24 horas.
func (r *BillingStoreCache) SetRates(ctx context.Context, rates map[string]float64) error {
	ratesMap := make(map[string]interface{}, len(rates))
	for k, v := range rates {
		ratesMap[k] = v
	}

	// Usamos una pipeline para asegurar que HSet y Expire se ejecuten en un solo viaje de red.
	pipe := r.rdb.Pipeline()
	pipe.HSet(ctx, billingdomain.LatestRatesKey, ratesMap)
	pipe.Expire(ctx, billingdomain.LatestRatesKey, billingdomain.LatestRatesExpiration)

	// Exec ejecuta los comandos en la pipeline.
	_, err := pipe.Exec(ctx)
	return err
}

// GetRates obtiene todas las tasas del Hash de Redis.
func (r *BillingStoreCache) GetRates(ctx context.Context) (map[string]float64, error) {
	// HGetAll devuelve un map[string]string
	stringRates, err := r.rdb.HGetAll(ctx, billingdomain.LatestRatesKey).Result()
	if err != nil {
		return nil, err
	}

	if len(stringRates) == 0 {
		// No es un error, pero no hay datos. Devolvemos redis.Nil
		// para que el llamador sepa que debe buscar en la API.
		return nil, redis.Nil
	}

	// Convertimos el map[string]string de vuelta a map[string]float64
	rates := make(map[string]float64, len(stringRates))
	for currency, rateStr := range stringRates {
		rateFloat, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			// Log de error, pero continuamos (podríamos querer fallar aquí)
			log.Printf("Error al parsear tasa para %s: %v", currency, err)
			continue
		}
		rates[currency] = rateFloat
	}

	return rates, nil
}

// DeleteRates elimina el hash completo de tasas.
func (r *BillingStoreCache) DeleteRates(ctx context.Context) error {
	return r.rdb.Del(ctx, billingdomain.LatestRatesKey).Err()
}

// SetSpecificTimeRateForCurrency guarda una tasa histórica con expiración.
func (r *BillingStoreCache) SetSpecificTimeRateForCurrency(ctx context.Context, date string, currency string, rate float64) error {
	// Creamos la clave dinámica: "rate:2025-11-15:VES"
	key := fmt.Sprintf("%s:%s:%s", billingdomain.RatePrefix, date, currency)

	// Usamos Set con el parámetro de expiración (rateExpiration)
	return r.rdb.Set(ctx, key, rate, billingdomain.RateExpiration).Err()
}

// GetSpecificTimeRateForCurrency obtiene una tasa histórica.
func (r *BillingStoreCache) GetSpecificTimeRateForCurrency(ctx context.Context, date string, currency string) (float64, error) {
	key := fmt.Sprintf("%s:%s:%s", billingdomain.RatePrefix, date, currency)

	rateStr, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		// Esto manejará redis.Nil si la clave no existe o ha expirado.
		return 0, err
	}

	// Convertimos el string de vuelta a float64
	return strconv.ParseFloat(rateStr, 64)
}
