package ratelimiter

import (
	"context"
	"log"
	"time"
)

type RateLimiter struct {
	config   *RateLimiterConfig
	database RateLimiterDatabaseInterface
}

type RateLimiterDatabaseInterface interface {
	BuscaTotalPorAPIKey(ctx context.Context, apikey string, chave int64) (int64, error)
	SomaRequisicaoPorAPIKey(ctx context.Context, apikey string, chave int64) error
	BuscaTotalPorIp(ctx context.Context, ip string, chave int64) (int64, error)
	SomaRequisicaoPorIp(ctx context.Context, ip string, chave int64) error
}

func NewRateLimiter(config *RateLimiterConfig, database RateLimiterDatabaseInterface) *RateLimiter {
	return &RateLimiter{
		config:   config,
		database: database,
	}
}

func (rl *RateLimiter) VerificaRegistraPorAPIKey(ctx context.Context, apikey string, datahoraRequest time.Time) (bool, error) {

	chave := datahoraRequest.Unix() / int64(rl.config.APIKeyDuration.Seconds())

	//Busca o Total por APIKey
	totalPorAPIKey, err := rl.database.BuscaTotalPorAPIKey(ctx, apikey, chave)

	if err != nil {
		return false, err
	}

	log.Printf("[APIKey] Saldo: %d", totalPorAPIKey)

	if totalPorAPIKey < int64(rl.config.IPMaxRequests) {
		err = rl.database.SomaRequisicaoPorAPIKey(ctx, apikey, chave)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (rl *RateLimiter) VerificaRegistraPorIp(ctx context.Context, ip string, datahoraRequest time.Time) (bool, error) {

	chave := datahoraRequest.Unix() / int64(rl.config.IPDuration.Seconds())

	//Buscar o Total por Ip
	totalPorIp, err := rl.database.BuscaTotalPorIp(ctx, ip, chave)

	if err != nil {
		return false, err
	}

	log.Printf("[IP] Saldo: %d", totalPorIp)

	if totalPorIp < int64(rl.config.APIKeyMaxRequests) {
		err = rl.database.SomaRequisicaoPorIp(ctx, ip, chave)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
