package domain

import (
	"context"

	"microservice/internal/pkg/errors"
	"microservice/models"
)

type DocumentDB interface {
	GetDocumentByID(ctx context.Context, id string) (models.Document, error)
	SaveDocument(ctx context.Context, doc models.Document) (string, error)
}
type Domain struct {
	db DocumentDB
}

func NewDomain(db DocumentDB) (*Domain, error) {
	return &Domain{
		db: db,
	}, nil
}

func (d *Domain) GetDocument(ctx context.Context, id string) (models.Document, error) {
	doc, err := d.db.GetDocumentByID(ctx, id)
	if err != nil {
		return models.Document{}, errors.Wrapf(err, "Failed to get document by id (%s) from DocumentDB", id)
	}

	return doc, nil
}

func (d *Domain) AddDocument(ctx context.Context, doc models.Document) (string, error) {
	id, err := d.db.SaveDocument(ctx, doc)
	if err != nil {
		return "", errors.Wrapf(err, "Failed save document (%v) in DocumentDB", doc)
	}

	return id, nil
}

func (d *Domain) TearDown(ctx context.Context) error {
	return nil
}
