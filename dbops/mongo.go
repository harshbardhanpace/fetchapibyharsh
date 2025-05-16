package dbops

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
}

func NewMongoRepository(client *mongo.Client, database string) MongoRepository {
	return &mongoRepository{
		client:   client,
		database: database,
	}
}

func (m *mongoRepository) GetMongoCollection(collectionName string) *mongo.Collection {
	collection := m.client.Database(m.database).Collection(collectionName)
	return collection
}

func (m *mongoRepository) GetMongoStatus() error {
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	err := m.client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (m *mongoRepository) InsertOne(collection string, document interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := m.client.Database(m.database).Collection(collection)
	_, err := coll.InsertOne(ctx, document)
	return err
}

func (m *mongoRepository) InsertMany(collection string, documents []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := m.client.Database(m.database).Collection(collection)
	_, err := coll.InsertMany(ctx, documents)
	return err
}

func (m *mongoRepository) FindOne(collection string, filter interface{}, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := m.client.Database(m.database).Collection(collection)
	return coll.FindOne(ctx, filter).Decode(result)
}

func (m *mongoRepository) Find(collection string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := m.client.Database(m.database).Collection(collection)
	return coll.Find(ctx, filter)
}

func (m *mongoRepository) FindOneAndCount(collectionName string, filter interface{}, result interface{}) (int64, error) {
	collection := m.GetMongoCollection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return 0, err
	}

	count, err := collection.CountDocuments(ctx, filter)
	return count, err
}

func (m *mongoRepository) FindDistinct(collectionName string, field string, filter interface{}) ([]interface{}, error) {
	collection := m.GetMongoCollection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	distinct, err := collection.Distinct(ctx, field, filter)
	if err != nil {
		return nil, err
	}

	return distinct, nil
}

func (m *mongoRepository) UpdateOne(collection string, filter, update interface{}, opts ...*options.UpdateOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := m.client.Database(m.database).Collection(collection)
	_, err := coll.UpdateOne(ctx, filter, update, opts...)
	return err
}

func (m *mongoRepository) DeleteOne(collection string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := m.client.Database(m.database).Collection(collection)
	delRes, err := coll.DeleteOne(ctx, filter, opts...)
	return delRes, err
}

func (m *mongoRepository) DeleteMany(collection string, filter interface{}, opts ...*options.DeleteOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := m.client.Database(m.database).Collection(collection)
	_, err := coll.DeleteMany(ctx, filter, opts...)
	return err
}

func (m *mongoRepository) ReplaceOne(collectionName string, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	collection := m.GetMongoCollection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.ReplaceOne(ctx, filter, replacement, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
