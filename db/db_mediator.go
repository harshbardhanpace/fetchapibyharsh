package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CallFindAllMongo = func(db MongoDatabase, collectionName string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return db.FindAllMongo(collectionName, filter, opts...)
}
