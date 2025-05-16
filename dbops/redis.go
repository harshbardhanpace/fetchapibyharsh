package dbops

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{client: client}
}

func (r *redisRepository) GetStatus() error {
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	if _, err := r.client.Ping(ctx).Result(); err != nil {
		return err
	}
	return nil
}

func (r *redisRepository) Set(key string, value interface{}, expiration time.Duration) error {
	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		strValue = string(jsonBytes)
	}
	return r.client.Set(context.Background(), key, strValue, expiration).Err()
}

func (r *redisRepository) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

func (r *redisRepository) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

func (r *redisRepository) Exists(key string) (int64, error) {
	return r.client.Exists(context.Background(), key).Result()
}

func (r *redisRepository) Increment(key string) int64 {
	return r.client.Incr(context.Background(), key).Val()
}

func (r *redisRepository) HSet(key, field string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.HSet(context.Background(), key, field, jsonValue).Err()
}

func (r *redisRepository) HGet(key, field string) (string, error) {
	return r.client.HGet(context.Background(), key, field).Result()
}

func (r *redisRepository) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(context.Background(), key).Result()
}

func (r *redisRepository) HDelete(key string, fields ...string) error {
	return r.client.HDel(context.Background(), key, fields...).Err()
}

func (r *redisRepository) HExists(key, field string) (bool, error) {
	return r.client.HExists(context.Background(), key, field).Result()
}

func (r *redisRepository) HKeys(key string) ([]string, error) {
	return r.client.HKeys(context.Background(), key).Result()
}
