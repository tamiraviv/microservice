package domain

import (
	"context"
	"reflect"
	"testing"

	"microservice/internal/pkg/errors"
	"microservice/mocks"
	"microservice/models"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestDomain_AddDocument(t *testing.T) {
	type dbAddDocumentMockData struct {
		times int
		err   error
	}

	successfulAddDocument := dbAddDocumentMockData{
		times: 1,
		err:   nil,
	}

	failedToAddDocument := dbAddDocumentMockData{
		times: 1,
		err:   errors.New("some-error"),
	}

	tests := []struct {
		name          string
		addDocumentMD dbAddDocumentMockData
		wantErr       bool
	}{
		{
			name:          "successful add document tp db expect no error",
			addDocumentMD: successfulAddDocument,
			wantErr:       false,
		},
		{
			name:          "failed to add document to db expect error",
			addDocumentMD: failedToAddDocument,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			id := uuid.New().String()

			db := mocks.NewMockDocumentDB(c)
			db.EXPECT().SaveDocument(gomock.Any(), gomock.AssignableToTypeOf(models.Document{})).Times(tt.addDocumentMD.times).Return(id, tt.addDocumentMD.err)

			d := &Domain{
				db: db,
			}

			docToAdd := models.Document{
				Name: "tamir",
				Doc: map[string]interface{}{
					"key": "value",
				},
			}

			got, err := d.AddDocument(context.TODO(), docToAdd)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, id) {
				t.Errorf("AddDocument() got = %v, want %v", got, id)
			}
		})
	}
}

func TestDomain_GetDocument(t *testing.T) {
	type dbGetDocumentMockData struct {
		times int
		err   error
	}

	successfulGetDocument := dbGetDocumentMockData{
		times: 1,
		err:   nil,
	}

	failedToGetDocument := dbGetDocumentMockData{
		times: 1,
		err:   errors.New("some-error"),
	}

	tests := []struct {
		name          string
		getDocumentMD dbGetDocumentMockData
		wantErr       bool
	}{
		{
			name:          "successful get document from db expect no error",
			getDocumentMD: successfulGetDocument,
			wantErr:       false,
		},
		{
			name:          "failed to get document from db expect error",
			getDocumentMD: failedToGetDocument,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			id := uuid.New().String()
			docToReturn := models.Document{
				Name: "tamir",
				Doc: map[string]interface{}{
					"key": "value",
				},
			}
			db := mocks.NewMockDocumentDB(c)
			db.EXPECT().GetDocumentByID(gomock.Any(), id, gomock.AssignableToTypeOf(&models.Document{})).
				Times(tt.getDocumentMD.times).
				Do(func(_ interface{}, _ interface{}, doc *models.Document) {
					*doc = docToReturn
				}).
				Return(tt.getDocumentMD.err)

			d := &Domain{
				db: db,
			}

			got, err := d.GetDocument(context.TODO(), id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, docToReturn) {
				t.Errorf("GetDocument() got = %v, want %v", got, docToReturn)
			}
		})
	}
}

func TestDomain_Teardown(t *testing.T) {
	type documentDBTearDownMockData struct {
		times int
		err   error
	}

	successfulTearDown := documentDBTearDownMockData{
		times: 1,
		err:   nil,
	}

	failedToTearDown := documentDBTearDownMockData{
		times: 1,
		err:   errors.New("some-error"),
	}

	tests := []struct {
		name       string
		TeardownMD documentDBTearDownMockData
		wantErr    bool
	}{
		{
			name:       "teardown successful expect no error",
			TeardownMD: successfulTearDown,
			wantErr:    false,
		},
		{
			name:       "failed to teardown document db expect error",
			TeardownMD: failedToTearDown,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			db := mocks.NewMockDocumentDB(c)
			db.EXPECT().Teardown(gomock.Any()).Times(tt.TeardownMD.times).Return(tt.TeardownMD.err)

			d := &Domain{
				db: db,
			}

			if err := d.Teardown(context.TODO()); (err != nil) != tt.wantErr {
				t.Errorf("Teardown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDomain(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "valid creation expect no error",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			db := mocks.NewMockDocumentDB(c)

			got, err := NewDomain(db)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			want := &Domain{db: db}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("NewDomain() got = %v, want %v", got, want)
			}
		})
	}
}
