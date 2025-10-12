package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RateLimiterDatabaseMock struct {
	mock.Mock
}

func (m *RateLimiterDatabaseMock) BuscaTotalPorAPIKey(ctx context.Context, apikey string, chave int64) (int64, error) {
	args := m.Called(chave)
	return args.Get(0).(int64), args.Error(1)
}

func (m *RateLimiterDatabaseMock) SomaRequisicaoPorAPIKey(ctx context.Context, apikey string, chave int64) error {
	args := m.Called(chave)
	return args.Error(0)
}

func (m *RateLimiterDatabaseMock) BuscaTotalPorIp(ctx context.Context, apikey string, chave int64) (int64, error) {
	args := m.Called(chave)
	return args.Get(0).(int64), args.Error(1)
}

func (m *RateLimiterDatabaseMock) SomaRequisicaoPorIp(ctx context.Context, apikey string, chave int64) error {
	args := m.Called(chave)
	return args.Error(0)
}

func Test_ChamadaPorAPIKey_Liberado(t *testing.T) {
	ctx := context.Background()

	rlDatabase := &RateLimiterDatabaseMock{}
	rlConfig := &RateLimiterConfig{
		APIKeyMaxRequests: 10,
		APIKeyDuration:    1 * time.Second,
		IPMaxRequests:     10,
		IPDuration:        1 * time.Second,
	}

	rl := NewRateLimiter(rlConfig, rlDatabase)

	chaveApiKey := int64(1760208325)
	dataResquest, err := time.Parse(time.RFC3339, "2025-10-11T18:45:25Z")

	rlDatabase.On("BuscaTotalPorAPIKey", chaveApiKey).Return(int64(0), nil)
	rlDatabase.On("SomaRequisicaoPorAPIKey", chaveApiKey).Return(nil)

	ok, err := rl.VerificaRegistra(ctx, "127.0.0.1", "ChaveXYZ", dataResquest)

	assert.Nil(t, err)
	assert.True(t, ok)
}

func Test_ChamadaPorAPIKey_Bloqueada_PorIP_Liberado(t *testing.T) {
	ctx := context.Background()

	rlDatabase := &RateLimiterDatabaseMock{}
	rlConfig := &RateLimiterConfig{
		APIKeyMaxRequests: 10,
		APIKeyDuration:    1 * time.Second,
		IPMaxRequests:     10,
		IPDuration:        1 * time.Second,
	}

	rl := NewRateLimiter(rlConfig, rlDatabase)

	chaveApiKey := int64(1760208325)
	dataResquest, err := time.Parse(time.RFC3339, "2025-10-11T18:45:25Z")

	rlDatabase.On("BuscaTotalPorAPIKey", chaveApiKey).Return(int64(10), nil)
	rlDatabase.On("SomaRequisicaoPorAPIKey", chaveApiKey).Return(nil)
	rlDatabase.On("BuscaTotalPorIp", chaveApiKey).Return(int64(0), nil)
	rlDatabase.On("SomaRequisicaoPorIp", chaveApiKey).Return(nil)

	ok, err := rl.VerificaRegistra(ctx, "127.0.0.1", "ChaveXYZ", dataResquest)

	assert.Nil(t, err)
	assert.True(t, ok)
}

func Test_ChamadaPorAPIKey_Bloqueada_PorIP_Bloqueado(t *testing.T) {
	ctx := context.Background()

	rlDatabase := &RateLimiterDatabaseMock{}
	rlConfig := &RateLimiterConfig{
		APIKeyMaxRequests: 10,
		APIKeyDuration:    1 * time.Second,
		IPMaxRequests:     10,
		IPDuration:        1 * time.Second,
	}

	rl := NewRateLimiter(rlConfig, rlDatabase)

	chaveApiKey := int64(1760208325)
	dataResquest, err := time.Parse(time.RFC3339, "2025-10-11T18:45:25Z")

	rlDatabase.On("BuscaTotalPorAPIKey", chaveApiKey).Return(int64(10), nil)
	rlDatabase.On("SomaRequisicaoPorAPIKey", chaveApiKey).Return(nil)
	rlDatabase.On("BuscaTotalPorIp", chaveApiKey).Return(int64(10), nil)

	ok, err := rl.VerificaRegistra(ctx, "127.0.0.1", "ChaveXYZ", dataResquest)

	assert.Nil(t, err)
	assert.False(t, ok)
}
