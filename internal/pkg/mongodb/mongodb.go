package mongodb

import (
	"context"
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

type Configuration interface {
	GetString(key string) (string, error)
	GetInt(key string) (int, error)
	GetDuration(key string) (time.Duration, error)
}

type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

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

	client.Database(database).Collection(collection)

	return &MongoDB{
		client:     client,
		collection: client.Database(database).Collection(collection),
	}, nil
}

func (m *MongoDB) GetDocumentByID(ctx context.Context, id string) (models.Document, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Document{}, errors.Errorf("id (%s) is not a valid ObjectID", id).SetType(errors.ErrorTypeBadRequest)
	}

	s := m.collection.FindOne(ctx, map[string]interface{}{"_id": objID})
	if err := s.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Document{}, errors.Errorf("Document with id (%s) was not found in mongodb", id).SetType(errors.ErrorTypeNotFound)
		}

		return models.Document{}, errors.Wrapf(err, "Failed to find document with id (%s) in mongodb", id)
	}

	var doc models.Document
	if err := s.Decode(&doc); err != nil {
		return models.Document{}, errors.Wrap(err, "Failed to decode document")
	}

	return doc, nil
}

func (m *MongoDB) SaveDocument(ctx context.Context, doc models.Document) (string, error) {
	res, err := m.collection.InsertOne(ctx, doc)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to insert document (%v) to mongodb", doc)
	}

	id := res.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}
