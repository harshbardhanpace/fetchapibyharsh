package db

import (
	"context"
	"net/url"
	"time"

	"space/constants"
	"space/loggerconfig"

	"space/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDbObj MongoDatabase

type MongoDb struct {
	client           *mongo.Client
	daoDb            *mongo.Database
	spaceDb          *mongo.Database
	contractSearchDb *mongo.Database
}

func (m *MongoDb) InitMongoClient(env string) error {
	defer models.HandlePanic()
	loggerconfig.Info("initiating mongo client for env= ", env+".mongo.mongo_uri")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoBaseURI := constants.MongoBase + url.QueryEscape(constants.MongoPass) + "@" + constants.MongoURI + "&" + constants.MongoQueryParam

	clientOpts := options.Client().ApplyURI(mongoBaseURI).SetMaxPoolSize(100)

	var err error
	m.client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		return err
	}
	// Ping the database to verify the connection and authentication
	err = m.client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	loggerconfig.Info("client vals=", m.client)

	m.spaceDb = m.client.Database(constants.ConfigSpaceDB)
	m.daoDb = m.client.Database(constants.ConfigDaoDB)
	m.contractSearchDb = m.client.Database(constants.ConfigContractSearchDB)

	if clientOpts.MaxPoolSize != nil {
		loggerconfig.Info("InitMongoClient Max pool size is:", *clientOpts.MaxPoolSize)
	} else {
		loggerconfig.Info("InitMongoClient Max pool size is not set explicitly.")
	}

	return err
}

func (m *MongoDb) GetMongoStatus() error {
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	err := m.client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func SetMongoDBObj(db MongoDatabase) {
	MongoDbObj = db
}

func GetMongoDBObj() MongoDatabase {
	return MongoDbObj
}

func (m *MongoDb) FindOneMongo(collectionName string, filter interface{}, result interface{}) error {
	collection := m.spaceDb.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(result)
	return err
}

func (m *MongoDb) FindManyMongo(collectionName string, filter interface{}, results interface{}) error {
	collection := m.spaceDb.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, results)
	return err
}

func (m *MongoDb) UpdateOneMongo(collectionName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	collection := m.spaceDb.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, filter, update, opts...)
	return err
}

func (m *MongoDb) UpdateOneMongoDao(collectionName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	collection := m.daoDb.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, filter, update, opts...)
	return err
}

func (m *MongoDb) FindAllMongo(collectionName string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	collection := m.spaceDb.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter, opts...)
	return cursor, err
}

func (m *MongoDb) InsertOneMongo(collectionName string, document interface{}) (*mongo.InsertOneResult, error) {
	collection := m.spaceDb.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, document)
	return result, err
}

func (m *MongoDb) DeleteOneMongo(collectionName string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	collection := m.spaceDb.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, filter, opts...)
	return result, err
}

func (m *MongoDb) DeleteMany(collectionName string, filter interface{}, opts ...*options.DeleteOptions) error {
	collection := m.spaceDb.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteMany(ctx, filter, opts...)
	return err
}

func (m *MongoDb) InsertMany(collectionName string, documents []interface{}, opts ...*options.InsertManyOptions) error {
	collection := m.spaceDb.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertMany(ctx, documents, opts...)
	return err
}

func (m *MongoDb) FindOneMongoDao(collectionName string, filter interface{}, result interface{}) error {
	collection := m.daoDb.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(result)
	return err
}

func (m *MongoDb) GetContractSearchDB() *mongo.Database {
	return m.contractSearchDb
}
