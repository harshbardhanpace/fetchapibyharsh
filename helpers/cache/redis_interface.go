package cache

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache interface {
	SetRedis(key string, val string, timeInMins time.Duration) error
	GetRedis(key string) *redis.StringCmd
	DelRedis(uemail string) *redis.IntCmd
	Exists(key string) *redis.IntCmd
	Incr(key string) *redis.IntCmd
	HSetField(key string, field string, value string) error
	HGetField(key string, field string) (string, error)
	Close() error
	InitRedis(string) error
	GetClientStatus() error
	GetOrderClientStatus() error
	FetchByScoreWithRange(key string, score float64, rangeSize int) ([]string, error)
	FetchByRoundedScoreWithRangeAndExpiry(key string, score float64, rangeSize int, expiry string) ([]string, error)
	SetWithTTL(key, value string, mins int) error
	HDel(key string, fields ...string) error
	SetRedisNoExpiry(key string, val string) error
	DeleteRedis(key string) error
	SAdd(key string, members ...string) error
	SRem(key string, members ...string) error
	LPush(key string, value ...string) error
}

type ContractCache interface {
	ContractCacheInit() error
	GetContractCacheStatus() error
	FlushCache() error
	AddToSortedSet(key string, score float64, val string)
	QuerySortedSetByPrefix(globalSearch, prefix string, startIndex, count int) ([]string, error)
	AddToHash(hash, key, value string) error
	AddToSet(key, value string, expiry time.Duration) error
	GetFromSet(key string) string
	GetFromHash(hash, key string) (error, string)
	GetAllFromHashWithPipeline(hash string, batchSize int64) (map[string]string, error)
}

type SmartCache interface {
	InitSmartCache() error
	GetStatusRedisSmartCache() error
	PerformNewSearch(exchange, searchTerm string, offset, capacity int, fuzzy bool) ([]string, error)
	ExecFTCommand(args []interface{}) ([]interface{}, error)
	GetFromHashSetNew(hash, key string) (string, error)
}
