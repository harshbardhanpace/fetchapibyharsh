package cache

import (
	"context"
	"space/constants"
	"space/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type ContractCacheRedisClient struct {
	contractCacheClient *redis.Client
}

var contractCacheClientObj ContractCache

func (c *ContractCacheRedisClient) ContractCacheInit() error {
	defer models.HandlePanic()

	c.contractCacheClient = redis.NewClient(&redis.Options{
		Addr:     constants.ContractCacheAddr,
		PoolSize: constants.ContractCachePoolSize,
		Password: constants.ContractCachePassword,
	})

	timeoutCtx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)

	_, tr := c.contractCacheClient.Ping(timeoutCtx).Result()
	if tr != nil {
		return tr
	}
	return nil
}

func (c *ContractCacheRedisClient) GetContractCacheStatus() error {
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	if _, err := c.contractCacheClient.Ping(ctx).Result(); err != nil {
		return err
	}

	return nil
}

func (c *ContractCacheRedisClient) FlushCache() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.contractCacheClient.FlushAll(timeoutCtx).Result()
	return err
}

func (c *ContractCacheRedisClient) AddToSortedSet(key string, score float64, val string) {
	member := &redis.Z{
		Score:  score,
		Member: val,
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c.contractCacheClient.ZAdd(timeoutCtx, key, *member)
}

func (c *ContractCacheRedisClient) QuerySortedSetByPrefix(globalSearch, prefix string, startIndex, count int) ([]string, error) {
	min := "[" + prefix
	max := "[" + prefix + "\xff"

	cmd := c.contractCacheClient.ZRangeByLex(context.Background(), globalSearch, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: int64(startIndex),
		Count:  int64(count),
	})

	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return cmd.Val(), nil
}

func (c *ContractCacheRedisClient) AddToHash(hash, key, value string) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.contractCacheClient.HSet(timeoutCtx, hash, key, value).Result()
	return err
}

func (c *ContractCacheRedisClient) AddToSet(key, value string, expiry time.Duration) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.contractCacheClient.Set(timeoutCtx, key, value, 0).Result()
	return err
}

func (c *ContractCacheRedisClient) GetFromSet(key string) string {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	val, _ := c.contractCacheClient.Get(timeoutCtx, key).Result()
	return val
}

func (c *ContractCacheRedisClient) GetFromHash(hash, key string) (error, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	val, err := c.contractCacheClient.HGet(ctx, hash, key).Result()
	return err, val
}

func (c *ContractCacheRedisClient) GetAllFromHashWithPipeline(hash string, batchSize int64) (map[string]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), constants.GetALLFromHashDelayTime*time.Second)
	defer cancel()

	result := make(map[string]string)
	var cursor uint64

	for {

		keys, newCursor, err := c.contractCacheClient.HScan(ctx, hash, cursor, "", batchSize).Result() // Fetch 1000 keys at a time
		if err != nil {
			return nil, err
		}

		for currentIndex := 0; currentIndex < len(keys); currentIndex += 2 {
			key := keys[currentIndex]
			value := keys[currentIndex+1]
			result[key] = value
		}

		cursor = newCursor

		if cursor == 0 {
			break
		}
	}

	return result, nil
}

func SetContractCacheClienttObj(contractCacheCliObj ContractCache) {
	contractCacheClientObj = contractCacheCliObj
}

func GetContractCacheClientObj() ContractCache {
	return contractCacheClientObj
}
