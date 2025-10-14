package main

import (
	"context"
	ratelimiter "fullcycle-desafio-tecnico-1/internal/rate_limiter"
	rateLimiterRedis "fullcycle-desafio-tecnico-1/internal/rate_limiter/redis"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var rl *ratelimiter.RateLimiter

func main() {
	// ctx := context.Background()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	rlConfig, err := ratelimiter.NewRateLimiterConfigFromEnvirontment()
	if err != nil {
		log.Fatal("Error trying to load rate limiter config")
		return
	}

	rlRedisConfig, err := rateLimiterRedis.NewRateLimiterRedisConfigFromEnvirontment()
	if err != nil {
		log.Fatal("Error trying to load rate limiter redis config")
		return
	}

	redisClient := rateLimiterRedis.NewRedisClient(rlRedisConfig)
	rlRedis := rateLimiterRedis.NewRateLimiterRedis(redisClient)
	rl = ratelimiter.NewRateLimiter(rlConfig, rlRedis)

	mux := IniciaServeMux()
	http.ListenAndServe(":8080", mux)
}

func IniciaServeMux() *http.ServeMux {

	mux := http.NewServeMux()
	mux.Handle("/", rateLimiterMiddleware(http.HandlerFunc(handler)))

	return mux
}

func rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		log.Println("[HTTP] Request received")

		apiKey := r.Header.Get("APIKEY")
		if apiKey == "" {
			log.Println("[HTTP] APIKey not found")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Não autorizado (Chave não encontrado)."))
			return
		}

		ip := extrairIp(r.RemoteAddr)
		//log.Println(fmt.Sprintf("[HTTP] Request received from IP: %s, Key: %s", ip, apiKey))
		log.Printf("[HTTP] Request received from IP: %s, Key: %s", ip, apiKey)

		tempo := time.Now()
		ok, err := rl.VerificaRegistra(ctx, ip, apiKey, tempo)

		if err != nil {
			log.Println("[HTTP] Error trying to register request")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Erro interno."))
			return
		}

		if !ok {
			log.Println("[HTTP] Rate limit exceeded")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Taxa de requisição excedida."))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sucesso."))
}

func extrairIp(RemoteAddr string) string {
	return strings.Split(RemoteAddr, ":")[0]
}
