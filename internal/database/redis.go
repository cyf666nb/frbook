package database

import (
	"context"
	"fmt"
	"time"

	"bookshare/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	return client, nil
}

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(cfg *config.RedisConfig) (*RedisClient, error) {
	client, err := NewRedis(cfg)
	if err != nil {
		return nil, err
	}
	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) GetJSON(ctx context.Context, key string, dest interface{}) error {
	return r.Get(ctx, key).Scan(dest)
}

func (r *RedisClient) Lock(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return r.SetNX(ctx, key, value, expiration).Result()
}

func (r *RedisClient) Unlock(ctx context.Context, key string, value string) error {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
	return r.Eval(ctx, script, []string{key}, value).Err()
}

func (r *RedisClient) IsLocked(ctx context.Context, key string) (bool, error) {
	exists, err := r.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *RedisClient) AddToDelayQueue(ctx context.Context, queueKey string, member string, score float64) error {
	return r.ZAdd(ctx, queueKey, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

func (r *RedisClient) GetExpiredFromDelayQueue(ctx context.Context, queueKey string, now float64) ([]string, error) {
	return r.ZRangeByScore(ctx, queueKey, &redis.ZRangeBy{
		Min:   "0",
		Max:   fmt.Sprintf("%f", now),
		Count: 100,
	}).Result()
}

func (r *RedisClient) RemoveFromDelayQueue(ctx context.Context, queueKey string, members ...string) error {
	values := make([]interface{}, len(members))
	for i, m := range members {
		values[i] = m
	}
	return r.ZRem(ctx, queueKey, values...).Err()
}
