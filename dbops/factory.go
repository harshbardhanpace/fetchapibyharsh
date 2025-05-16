package dbops

import (
	"context"
	"database/sql"
	"fmt"
	"space/constants"
	"space/loggerconfig"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseType string

const (
	RedisType    = "redis"
	MongoType    = "mongo"
	PostgresType = "postgres"
)

type DatabaseConfig struct {
	Type               DatabaseType
	Host               string
	Port               string
	Username           string
	Password           string
	Database           string
	PoolSize           int
	MaxConnection      int
	MaxIdleConnections int
}

type DatabaseConnector interface {
	Connect() error
	Close() error
	GetRepository() interface{}
}

type RedisConnector struct {
	client *redis.Client
	config DatabaseConfig
}

func (r *RedisConnector) Connect() error {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.config.Host, r.config.Port),
		Username: r.config.Username,
		Password: r.config.Password,
		PoolSize: r.config.PoolSize,
	}
	r.client = redis.NewClient(opts)
	_, err := r.client.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	logrus.Printf("Connected to Redis at %s:%s", r.config.Host, r.config.Port)
	return nil
}

func (r *RedisConnector) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

func (r *RedisConnector) GetRepository() interface{} {
	return NewRedisRepository(r.client)
}

type MongoConnector struct {
	client *mongo.Client
	config DatabaseConfig
}

func (m *MongoConnector) Connect() error {
	mongoURI := m.config.Host + m.config.Password + "@" + m.config.Port + "&" + constants.MongoQueryParam
	clientOptions := options.Client().ApplyURI(mongoURI).SetMaxPoolSize(20)
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	m.client = client
	logrus.Printf("Connected to MongoDB at: %s", mongoURI)
	return nil
}

func (m *MongoConnector) Close() error {
	if m.client != nil {
		return m.client.Disconnect(context.Background())
	}
	return nil
}

func (m *MongoConnector) GetRepository() interface{} {
	return NewMongoRepository(m.client, m.config.Database)
}

type PostgresConnector struct {
	db     *sql.DB
	config DatabaseConfig
}

func (p *PostgresConnector) Connect() error {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s TimeZone=Asia/Kolkata sslrootcert=%s",
		p.config.Host, p.config.Port, p.config.Username,
		p.config.Password, p.config.Database, constants.CertificateEnabled, constants.CertificatePath)

	db, err := sql.Open("postgres", connectionString+" connect_timeout=1")
	if err != nil {
		return fmt.Errorf("failed to connect to Postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping Postgres: %w", err)
	}
	db.SetMaxOpenConns(p.config.MaxConnection)
	db.SetMaxIdleConns(p.config.MaxIdleConnections)
	db.SetConnMaxLifetime(24 * time.Hour)
	p.db = db
	err = db.PingContext(context.Background())
	if err != nil {
		loggerconfig.Error(" Error ping : " + err.Error())
		return err
	}
	logrus.Printf("Connected to Postgres at %s:%s", p.config.Host, p.config.Port)
	return nil
}

func (p *PostgresConnector) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PostgresConnector) GetRepository() interface{} {
	return NewPostgresRepository(p.db)
}

func (p *PostgresConnector) GetDB() *sql.DB {
	return p.db
}

type DatabaseConnectorFactory struct{}

func NewDatabaseConnectorFactory() *DatabaseConnectorFactory {
	return &DatabaseConnectorFactory{}
}

func (f *DatabaseConnectorFactory) GetConnector(config DatabaseConfig) (DatabaseConnector, error) {
	var connector DatabaseConnector

	switch config.Type {
	case RedisType:
		connector = &RedisConnector{config: config}
	case MongoType:
		connector = &MongoConnector{config: config}
	case PostgresType:
		connector = &PostgresConnector{config: config}
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}

	err := connector.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", config.Type, err)
	}

	return connector, nil
}
