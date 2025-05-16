package dbops

import (
	"database/sql"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RedisRepository interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (int64, error)
	Increment(key string) int64

	HSet(key, field string, value interface{}) error
	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	HDelete(key string, fields ...string) error
	HExists(key, field string) (bool, error)
	HKeys(key string) ([]string, error)

	GetStatus() error
}

type MongoRepository interface {
	GetMongoCollection(collectionName string) *mongo.Collection
	GetMongoStatus() error

	InsertOne(collection string, document interface{}) error
	InsertMany(collection string, documents []interface{}) error
	FindOne(collection string, filter interface{}, result interface{}) error
	Find(collection string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOneAndCount(collectionName string, filter interface{}, result interface{}) (int64, error)
	FindDistinct(collectionName string, field string, filter interface{}) ([]interface{}, error)
	UpdateOne(collection string, filter, update interface{}, opts ...*options.UpdateOptions) error
	DeleteOne(collection string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(collection string, filter interface{}, opts ...*options.DeleteOptions) error
	ReplaceOne(collectionName string, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
}

type PostgresRepository interface {
	GetStatus() error

	Insert(query string, args ...any) (*sql.Rows, error)
	Update(query string, args ...any) (*sql.Rows, error)
	Fetch(query string, args ...any) (*sql.Rows, error)
	Delete(query string, args ...any) (*sql.Rows, error)
}

// DatabaseFactory is the abstract factory interface
type DatabaseFactory interface {
	CreateRedisRepository() (RedisRepository, error)
	CreateMongoRepository() (MongoRepository, error)
	CreatePostgresRepository() (PostgresRepository, error)
}
