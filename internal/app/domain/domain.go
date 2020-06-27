package domain

import (
	"context"

	"microservice/internal/pkg/errors"
	"microservice/models"
)

// DocumentDB expose CRUD related operations for document
type DocumentDB interface {
	GetDocumentByID(ctx context.Context, id string, result interface{}) error
	SaveDocument(ctx context.Context, doc models.Document) (string, error)
	Teardown(ctx context.Context) error
}

// Domain implement a Domain Service
type Domain struct {
	db DocumentDB
}

// NewDomain returns a new instance of the Domain struct
func NewDomain(db DocumentDB) (*Domain, error) {
	return &Domain{
		db: db,
	}, nil
}

// GetDocument gets an id and return the document of that id
func (d *Domain) GetDocument(ctx context.Context, id string) (models.Document, error) {
	var doc models.Document
	if err := d.db.GetDocumentByID(ctx, id, &doc); err != nil {
		return models.Document{}, errors.Wrapf(err, "Failed to get document by id (%s) from DocumentDB", id)
	}

	return doc, nil
}

// AddDocument gets a document, save it to the document db and return id of that document for further queries
func (d *Domain) AddDocument(ctx context.Context, doc models.Document) (string, error) {
	id, err := d.db.SaveDocument(ctx, doc)
	if err != nil {
		return "", errors.Wrapf(err, "Failed save document (%v) in DocumentDB", doc)
	}

	return id, nil
}

// Teardown closes every open connection of the domain
func (d *Domain) Teardown(ctx context.Context) error {
	if err := d.db.Teardown(ctx); err != nil {
		return errors.Wrap(err, "Failed to gracefully teardown document db")
	}
	return nil
}
