package db

//go:generate mockgen -source=./connection.go -package=mocks -destination=../../../mocks/mock_connection.go

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// IConnection is a low level interface to the database.
type IConnection interface {
	// CleanUpDatabase closes the connection to the database. You should defer a call to this method after calling
	// NewDbConnection. Any errors are unexpected.
	CleanUpDatabase() error

	// Insert inserts a document into the database at the specified collection. The document is inserted with a new
	// unique ID, or the _id field is used if it is already set. The result of the insert is returned.
	// This method returns an error if a document with the given ID already exists. Use mongo.IsDuplicateKeyError
	// to check for this error. Any other errors are unexpected.
	Insert(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error)

	// GetIDs returns the IDs of all documents in the specified collection. The IDs are returned as a slice of
	// bson.M objects, where each object has a single key "_id" with the ID as the value.
	// If the collection does not exist, an empty slice is returned. Any errors are unexpected.
	GetIDs(ctx context.Context, collection string, resultIds *[]bson.M) error

	// FindOne returns a single document from the specified collection that matches the given filter. The filter
	// should be a bson.M object. If no document is found, calling the Decode method on the returned SingleResult
	// will return a mongo.ErrNoDocuments error.
	FindOne(ctx context.Context, collection string, filter interface{}) *mongo.SingleResult

	// DeleteOne deletes a single document from the specified collection that matches the given filter. The filter
	// should be a bson.M object. The result of the delete operation is returned. If no matching document is found, the
	// DeletedCount field of the result will be 0, and no error will be returned.
	DeleteOne(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error)

	// DropCollection drops a given collection. This is a destructive operation and should only be used for testing.
	DropCollection(ctx context.Context, collection string) error
}

type connection struct {
	database *mongo.Database
	client   *mongo.Client
}

// NewDbConnection creates a new connection to the database. You should defer a call to the CleanUpDatabase method
// on the returned IConnection object.
// Any errors are unexpected.
func NewDbConnection(config *Config) (IConnection, error) {
	m := connection{}
	return &m, m.setupDatabase(config)
}

func (m *connection) setupDatabase(config *Config) error {
	opts := options.Client()
	opts.ApplyURI("mongodb://" + config.User + ":" + config.Password + "@" + config.Host + ":" + "27017" + "/" + config.Db)
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	m.client, err = mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	// ensure connection was successful
	if err = m.client.Ping(context.Background(), nil); err != nil {
		return err
	}

	m.database = m.client.Database(config.Db, options.Database())
	return nil
}

func (m *connection) CleanUpDatabase() error {
	return m.client.Disconnect(context.Background())
}

func (m *connection) Insert(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	return m.database.Collection(collection).InsertOne(ctx, document)
}

func (m *connection) GetIDs(ctx context.Context, collection string, resultIds *[]bson.M) error {
	// we are only interested in the _id field
	opts := options.Find().SetProjection(bson.D{{"_id", 1}})

	cursor, err := m.database.Collection(collection).Find(ctx, bson.D{}, opts)
	if err != nil {
		return err
	}

	return cursor.All(ctx, resultIds)
}

func (m *connection) FindOne(ctx context.Context, collection string, filter interface{}) *mongo.SingleResult {
	return m.database.Collection(collection).FindOne(ctx, filter)
}

func (m *connection) DeleteOne(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	return m.database.Collection(collection).DeleteOne(ctx, filter)
}

func (m *connection) DropCollection(ctx context.Context, collection string) error {
	return m.database.Collection(collection).Drop(ctx)
}
