package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr      string
	Password  string
	DB        int
	KeyPrefix string
}

type JSONCache struct {
	client    *redis.Client
	keyPrefix string
}

func NewJSONCache(cfg RedisConfig) *JSONCache {
	addr := strings.TrimSpace(cfg.Addr)
	if addr == "" {
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &JSONCache{
		client:    client,
		keyPrefix: strings.Trim(strings.TrimSpace(cfg.KeyPrefix), ":"),
	}
}

func (c *JSONCache) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

func GetOrSet[T any](
	ctx context.Context,
	cache *JSONCache,
	key string,
	ttl time.Duration,
	fetch func(context.Context) (T, error),
) (T, error) {
	if cache == nil || cache.client == nil {
		log.Printf("dictionary cache disabled key=%s source=postgres", key)
		return fetchFromPostgres(ctx, key, fetch)
	}

	cacheKey := cache.key(key)
	raw, err := cache.client.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var value T
		if jsonErr := json.Unmarshal(raw, &value); jsonErr == nil {
			log.Printf("dictionary cache hit key=%s source=redis", cacheKey)
			return value, nil
		}
		log.Printf("dictionary cache invalid_json key=%s source=postgres", cacheKey)
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		log.Printf("dictionary cache get_error key=%s source=postgres error=%v", cacheKey, err)
		return fetchFromPostgres(ctx, cacheKey, fetch)
	}

	if errors.Is(err, redis.Nil) {
		log.Printf("dictionary cache miss key=%s source=postgres", cacheKey)
	}

	value, fetchErr := fetchFromPostgres(ctx, cacheKey, fetch)
	if fetchErr != nil {
		return value, fetchErr
	}

	payload, marshalErr := json.Marshal(value)
	if marshalErr != nil {
		log.Printf("dictionary cache marshal_error key=%s error=%v", cacheKey, marshalErr)
		return value, nil
	}

	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	if setErr := cache.client.Set(ctx, cacheKey, payload, ttl).Err(); setErr != nil {
		log.Printf("dictionary cache set_error key=%s error=%v", cacheKey, setErr)
		return value, nil
	}
	log.Printf("dictionary cache set key=%s ttl=%s", cacheKey, ttl)

	return value, nil
}

func fetchFromPostgres[T any](ctx context.Context, key string, fetch func(context.Context) (T, error)) (T, error) {
	value, err := fetch(ctx)
	if err != nil {
		log.Printf("dictionary cache postgres_error key=%s error=%v", key, err)
		return value, err
	}
	log.Printf("dictionary cache postgres_ok key=%s", key)
	return value, nil
}

func (c *JSONCache) key(key string) string {
	key = strings.Trim(key, ":")
	if c.keyPrefix == "" {
		return key
	}
	return c.keyPrefix + ":" + key
}
