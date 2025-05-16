package cache

import (
	"context"
	"encoding/json"
	"math"
	"space/constants"
	"space/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client      *redis.Client
	OrderClient *redis.Client
}

var redisClientObj RedisCache

func NewRedisClient(url string) (*redis.Client, error) {
	defer models.HandlePanic()
	opts := &redis.Options{
		Addr: url,
	}

	client := redis.NewClient(opts)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (r *RedisClient) SetRedis(key string, val string, timeInMins time.Duration) error {
	return r.Client.Set(context.Background(), key, val, time.Minute*timeInMins).Err()
}

func (r *RedisClient) GetRedis(key string) *redis.StringCmd {
	return r.Client.Get(context.Background(), key)
}

func (r *RedisClient) DelRedis(uemail string) *redis.IntCmd {
	return r.Client.Del(context.Background(), uemail)
}

func (r *RedisClient) Exists(key string) *redis.IntCmd {
	return r.Client.Exists(context.Background(), key)
}

func (r *RedisClient) Incr(key string) *redis.IntCmd {
	return r.Client.Incr(context.Background(), key)
}

func (r *RedisClient) HSetField(key string, field string, value string) error {
	return r.OrderClient.HSet(context.Background(), key, field, value).Err()
}

func (r *RedisClient) HGetField(key string, field string) (string, error) {
	return r.OrderClient.HGet(context.Background(), key, field).Result()
}
func (r *RedisClient) HDel(key string, fields ...string) error {
	return r.Client.HDel(context.Background(), key, fields...).Err()
}

func (r *RedisClient) SetRedisNoExpiry(key string, val string) error {
	return r.Client.Set(context.Background(), key, val, 0).Err() // 0 means no expiry
}

func (r *RedisClient) DeleteRedis(key string) error {
	return r.Client.Del(context.Background(), key).Err()
}

func (r *RedisClient) SAdd(key string, members ...string) error {
	// Key is name of set, value is member of set
	return r.Client.SAdd(context.Background(), key, members).Err()
}

func (r *RedisClient) SRem(key string, members ...string) error {
	// Key is the name of the set, members are the elements to be removed
	return r.Client.SRem(context.Background(), key, members).Err()
}

func (r *RedisClient) LPush(key string, values ...string) error {
	// Push values to the left of the list
	return r.Client.LPush(context.Background(), key, values).Err()
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}

func (r *RedisClient) InitRedis(redisVal string) error {
	if redisVal == constants.Redis || redisVal == constants.MainRedis {
		client, err := NewRedisClient(constants.RedisUrl)
		if err != nil {
			return err
		}
		r.Client = client
	}

	if redisVal == constants.Redis || redisVal == constants.OrderRedis {
		orderClient, err := NewRedisClient(constants.OrderRedisUrl)
		if err != nil {
			return err
		}
		r.OrderClient = orderClient
	}

	return nil
}

func (r *RedisClient) GetClientStatus() error {
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	if _, err := r.Client.Ping(ctx).Result(); err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) GetOrderClientStatus() error {
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	if _, err := r.OrderClient.Ping(ctx).Result(); err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) FetchByScoreWithRange(key string, score float64, rangeSize int) ([]string, error) {
	fixedDifference := 500.0
	halfRangeSize := rangeSize / 2

	min := &redis.ZRangeBy{
		Min: strconv.FormatFloat(score-fixedDifference, 'f', -1, 64),
		Max: strconv.FormatFloat(score+fixedDifference, 'f', -1, 64),
	}

	result, err := r.Client.ZRangeByScore(context.Background(), key, min).Result()
	if err != nil {
		return nil, err
	}

	var filteredResultLen []string
	var expiry string
	set := false

	for _, member := range result {
		var data models.OptionData
		if err := json.Unmarshal([]byte(member), &data); err == nil {
			if !set {
				expiry = data.ExpiryRaw
				set = true
			}
			if data.ExpiryRaw == expiry {
				filteredResultLen = append(filteredResultLen, member)
			}
		}
	}

	interval := fixedDifference * 2 / float64(len(filteredResultLen)-1)
	closestIntScore := math.Round(score/interval) * interval

	min = &redis.ZRangeBy{
		Min: strconv.FormatFloat(closestIntScore-float64(halfRangeSize)*interval, 'f', -1, 64),
		Max: strconv.FormatFloat(closestIntScore+float64(halfRangeSize)*interval, 'f', -1, 64),
	}

	result, err = r.Client.ZRangeByScore(context.Background(), key, min).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *RedisClient) FetchByRoundedScoreWithRangeAndExpiry(key string, score float64, rangeSize int, expiry string) ([]string, error) {
	fixedDifference := 500.0
	halfRangeSize := rangeSize / 2

	min := &redis.ZRangeBy{
		Min: strconv.FormatFloat(score-fixedDifference, 'f', -1, 64),
		Max: strconv.FormatFloat(score+fixedDifference, 'f', -1, 64),
	}

	result, err := r.Client.ZRangeByScore(context.Background(), key, min).Result()
	if err != nil {
		return nil, err
	}

	var filteredResultLen []string
	for _, member := range result {
		var data models.OptionData
		if err := json.Unmarshal([]byte(member), &data); err == nil {
			if data.ExpiryRaw == expiry {
				filteredResultLen = append(filteredResultLen, member)
			}
		}
	}

	interval := fixedDifference * 2 / float64(len(filteredResultLen)-1)
	closestIntScore := math.Round(score/interval) * interval

	min.Min = strconv.FormatFloat(closestIntScore-float64(halfRangeSize)*interval, 'f', -1, 64)
	min.Max = strconv.FormatFloat(closestIntScore+float64(halfRangeSize)*interval, 'f', -1, 64)

	result, err = r.Client.ZRangeByScore(context.Background(), key, min).Result()
	if err != nil {
		return nil, err
	}

	var filteredResult []string
	for _, member := range result {
		var data models.OptionData
		if err := json.Unmarshal([]byte(member), &data); err == nil {
			if data.ExpiryRaw == expiry {
				filteredResult = append(filteredResult, member)
			}
		}
	}

	return filteredResult, nil
}

func (r *RedisClient) SetWithTTL(key, value string, mins int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.Client.Set(ctx, key, value, time.Duration(time.Minute)*time.Duration(mins)).Err()
}

func SetRedisClientObj(redisCliObj RedisCache) {
	redisClientObj = redisCliObj
}

func GetRedisClientObj() RedisCache {
	return redisClientObj
}
