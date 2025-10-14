package rateLimiterRedis

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RateLimiterRedis struct {
	client *redis.Client
}

func NewRateLimiterRedis(client *redis.Client) *RateLimiterRedis {
	return &RateLimiterRedis{
		client: client,
	}
}

func NewRedisClient(config *RateLimiterRedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password, // No password set
		DB:       0,               // Use default DB
		Protocol: 2,               // Connection protocol
	})
}

func (rl *RateLimiterRedis) BuscaTotalPorAPIKey(ctx context.Context, apikey string, chave int64) (int64, error) {
	id := retornaAPIKeyId(apikey, chave)
	valor, err := rl.client.Get(ctx, id).Result()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return -1, err
	}

	valorInt, err := strconv.ParseInt(valor, 10, 64)
	if err != nil {
		return -1, err
	}

	return valorInt, nil
}

func (rl *RateLimiterRedis) SomaRequisicaoPorAPIKey(ctx context.Context, apikey string, chave int64) error {
	id := retornaAPIKeyId(apikey, chave)
	err := rl.client.Incr(ctx, id).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rl *RateLimiterRedis) BuscaTotalPorIp(ctx context.Context, ip string, chave int64) (int64, error) {
	id := retornaIpId(ip, chave)
	valor, err := rl.client.Get(ctx, id).Result()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return -1, err
	}

	valorInt, err := strconv.ParseInt(valor, 10, 64)
	if err != nil {
		return -1, err
	}

	return valorInt, nil
}

func (rl *RateLimiterRedis) SomaRequisicaoPorIp(ctx context.Context, ip string, chave int64) error {
	id := retornaIpId(ip, chave)
	err := rl.client.Incr(ctx, id).Err()

	if err != nil {
		return err
	}

	return nil
}

func retornaAPIKeyId(apikey string, chave int64) string {
	return "KEY-" + apikey + "-" + strconv.FormatInt(chave, 10)
}

func retornaIpId(ip string, chave int64) string {
	return "IP-" + ip + "-" + strconv.FormatInt(chave, 10)
}
