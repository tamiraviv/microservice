package mongodb

import (
	"context"
	"reflect"
	"strings"
	"time"

	"microservice/internal/pkg/errors"
	"microservice/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoBaseKey       = "mongo"
	mongoHostsKey      = mongoBaseKey + ".hosts"
	mongoUsernameKey   = mongoBaseKey + ".username"
	mongoPasswordKey   = mongoBaseKey + ".password"
	mongoDatabaseKey   = mongoBaseKey + ".database"
	mongoCollectionKey = mongoBaseKey + ".collection"
)

// Configuration expose an interface of configuration related actions
type Configuration interface {
	GetString(key string) (string, error)
	GetInt(key string) (int, error)
	GetDuration(key string) (time.Duration, error)
}

// MongoDB client fpr mongodb which specifies which database and collection to use
type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewClient returns a new instance of the MongoDB struct
func NewClient(ctx context.Context, conf Configuration) (*MongoDB, error) {
	o := options.Client()

	hosts, err := conf.GetString(mongoHostsKey)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to get mongo hosts from configuration key (%s)", mongoHostsKey)
	}

	username, err := conf.GetString(mongoUsernameKey)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to get mongo username from configuration key (%s)", mongoUsernameKey)
	}

	password, err := conf.GetString(mongoPasswordKey)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to get mongo password from configuration key (%s)", mongoPasswordKey)
	}

	database, err := conf.GetString(mongoDatabaseKey)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to get mongo database from configuration key (%s)", mongoDatabaseKey)
	}

	collection, err := conf.GetString(mongoCollectionKey)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to get mongo collection from configuration key (%s)", mongoCollectionKey)
	}

	o.SetHosts(strings.Split(hosts, ","))
	o.SetAuth(options.Credential{
		Username: username,
		Password: password,
	})

	client, err := mongo.Connect(ctx, o)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to connect to mongodb")
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, errors.Wrap(err, "Fail to ping to mongodb server")
	}

	return &MongoDB{
		client:     client,
		collection: client.Database(database).Collection(collection),
	}, nil
}

// GetDocumentByID get document by ID from mongodb, and put it in the parameter 'result'.
// Note that result should be a pointer the the desired type
func (m *MongoDB) GetDocumentByID(ctx context.Context, id string, result interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Errorf("id (%s) is not a valid ObjectID", id).SetType(errors.ErrorTypeBadRequest)
	}

	s := m.collection.FindOne(ctx, map[string]interface{}{"_id": objID})
	if err := s.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.Errorf("Document with id (%s) was not found in mongodb", id).SetType(errors.ErrorTypeNotFound)
		}

		return errors.Wrapf(err, "Failed to find document with id (%s) in mongodb", id).SetType(errors.ErrorTypeInternal)
	}

	if err := s.Decode(result); err != nil {
		return errors.Wrapf(err, "Failed to decode document to result type (%s)", reflect.TypeOf(result)).SetType(errors.ErrorTypeBadRequest)
	}

	return nil
}

// SaveDocument add document to mongodb, return the id of the document
func (m *MongoDB) SaveDocument(ctx context.Context, doc models.Document) (string, error) {
	res, err := m.collection.InsertOne(ctx, doc)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to insert document (%v) to mongodb", doc).SetType(errors.ErrorTypeInternal)
	}

	id := res.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}

// Teardown disconnect from mongodb client
func (m *MongoDB) Teardown(ctx context.Context) error {
	if err := m.client.Disconnect(ctx); err != nil {
		return errors.Wrap(err, "Failed to close mongodb connections")
	}

	return nil
}
