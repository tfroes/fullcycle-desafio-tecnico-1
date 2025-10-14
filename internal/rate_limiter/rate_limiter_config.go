package ratelimiter

import (
	"os"
	"strconv"
	"time"
)

type RateLimiterConfig struct {
	// IP Rate Limiter
	IPMaxRequests int64
	IPDuration    time.Duration

	// API Key Rate Limiter
	APIKeyMaxRequests int64
	APIKeyDuration    time.Duration
}

func NewRateLimiterConfig(
	ipMaxRequests int64,
	ipInterval string,
	apiKeyMaxRequests int64,
	apiKeyInterval string) (*RateLimiterConfig, error) {

	ipDuration, err := time.ParseDuration(ipInterval)
	if err != nil {
		return nil, err
	}

	apiKeyDuration, err := time.ParseDuration(apiKeyInterval)
	if err != nil {
		return nil, err
	}

	return &RateLimiterConfig{
		IPMaxRequests:     ipMaxRequests,
		IPDuration:        ipDuration,
		APIKeyMaxRequests: apiKeyMaxRequests,
		APIKeyDuration:    apiKeyDuration,
	}, nil
}

func NewRateLimiterConfigFromEnvirontment() (*RateLimiterConfig, error) {

	ipMaxRequestsString := os.Getenv("IP_MAX_REQUESTS")
	ipInterval := os.Getenv("IP_INTERVAL")
	apiKeyMaxRequestsString := os.Getenv("APIKEY_MAX_REQUESTS")
	apiKeyInterval := os.Getenv("APIKEY_INTERVAL")

	ipMaxRequests, err := strconv.ParseInt(ipMaxRequestsString, 10, 64)
	if err != nil {
		return nil, err
	}

	apiKeyMaxRequests, err := strconv.ParseInt(apiKeyMaxRequestsString, 10, 64)
	if err != nil {
		return nil, err
	}

	ipDuration, err := time.ParseDuration(ipInterval)
	if err != nil {
		return nil, err
	}

	apiKeyDuration, err := time.ParseDuration(apiKeyInterval)
	if err != nil {
		return nil, err
	}

	return &RateLimiterConfig{
		IPMaxRequests:     ipMaxRequests,
		IPDuration:        ipDuration,
		APIKeyMaxRequests: apiKeyMaxRequests,
		APIKeyDuration:    apiKeyDuration,
	}, nil
}
