package rateLimiterRedis

import (
	"os"
)

type RateLimiterRedisConfig struct {
	Address  string
	Password string
}

func NewRateLimiterRedisConfig(
	address string,
	password string) (*RateLimiterRedisConfig, error) {

	return &RateLimiterRedisConfig{
		Address:  address,
		Password: password,
	}, nil
}

func NewRateLimiterRedisConfigFromEnvirontment() (*RateLimiterRedisConfig, error) {

	address := os.Getenv("RL_REDIS_ADDRESS")
	password := os.Getenv("RL_REDIS_PASSWORD")

	return &RateLimiterRedisConfig{
		Address:  address,
		Password: password,
	}, nil
}
