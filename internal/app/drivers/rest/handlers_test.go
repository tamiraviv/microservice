package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"microservice/internal/pkg/errors"
	"microservice/models"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"microservice/mocks"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAdapter_getDocument(t *testing.T) {
	type domainServiceGetDocumentMockData struct {
		times int
		err   error
		doc   models.Document
	}

	validDoc := models.Document{}

	unmarshallabledDoc := models.Document{
		Doc: map[string]interface{}{
			"invalid-doc": make(chan int),
		},
	}

	successfulGetValidDocument := domainServiceGetDocumentMockData{
		times: 1,
		err:   nil,
		doc:   validDoc,
	}

	badRequest := domainServiceGetDocumentMockData{
		times: 1,
		err:   errors.New("bad-request").SetType(errors.ErrorTypeBadRequest),
	}

	documentNotFound := domainServiceGetDocumentMockData{
		times: 1,
		err:   errors.New("not-found").SetType(errors.ErrorTypeNotFound),
	}

	failedToGetDocument := domainServiceGetDocumentMockData{
		times: 1,
		err:   errors.New("some-error").SetType(errors.ErrorTypeInternal),
	}

	GetInvalidDocument := domainServiceGetDocumentMockData{
		times: 1,
		err:   nil,
		doc:   unmarshallabledDoc,
	}

	tests := []struct {
		name                       string
		domainServiceGetDocumentMD domainServiceGetDocumentMockData
		wantedStatusCode           int
		wantErr                    bool
	}{
		{
			name:                       "get document successfully expect status OK (200)",
			domainServiceGetDocumentMD: successfulGetValidDocument,
			wantedStatusCode:           http.StatusOK,
			wantErr:                    false,
		},
		{
			name:                       "get bad id expect status bad request (400)",
			domainServiceGetDocumentMD: badRequest,
			wantedStatusCode:           http.StatusBadRequest,
			wantErr:                    true,
		},
		{
			name:                       "id doesn't exist in db expect status not found (404)",
			domainServiceGetDocumentMD: documentNotFound,
			wantedStatusCode:           http.StatusNotFound,
			wantErr:                    true,
		},
		{
			name:                       "failed to get document from db expect status internal server error (500)",
			domainServiceGetDocumentMD: failedToGetDocument,
			wantedStatusCode:           http.StatusInternalServerError,
			wantErr:                    true,
		},
		{
			name:                       "failed to marshal response document from db expect status internal server error (500)",
			domainServiceGetDocumentMD: GetInvalidDocument,
			wantedStatusCode:           http.StatusInternalServerError,
			wantErr:                    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			id := uuid.New().String()

			domainService := mocks.NewMockDomainService(c)
			domainService.EXPECT().GetDocument(gomock.Any(), id).
				Times(tt.domainServiceGetDocumentMD.times).
				Return(tt.domainServiceGetDocumentMD.doc, tt.domainServiceGetDocumentMD.err)

			s := &Adapter{
				domainSvc: domainService,
			}

			r := chi.NewRouter()
			r.Route("/documents", func(r chi.Router) {
				r.Get("/{id}", s.getDocument)
			})

			ts := httptest.NewServer(r)
			defer ts.Close()

			res, body := testRequest(t, ts, http.MethodGet, fmt.Sprintf("/documents/%s", id), nil)
			statusCodeCheck(t, res, tt.wantedStatusCode)

			if tt.wantErr {
				return
			}

			var respDoc models.Document
			if err := json.Unmarshal(body, &respDoc); err != nil {
				t.Fatalf("Failed to unmarshal response body to 'Document'. Error: %s", err)
			}

			if !reflect.DeepEqual(respDoc, tt.domainServiceGetDocumentMD.doc) {
				t.Fatalf("getDocument() got = %v, want %v", respDoc, tt.domainServiceGetDocumentMD.doc)
			}
		})
	}
}

func TestAdapter_addDocument(t *testing.T) {
	type jsonSchemaValidatorMockData struct {
		times int
		err   error
	}

	type domainServiceAddDocumentMockData struct {
		times int
		err   error
	}

	validDoc := models.Document{
		Name: "tamir",
		Doc: map[string]interface{}{
			"lastName": "Aviv",
		},
	}

	invalidDoc := models.Document{
		Name: "empty-document",
	}

	validJSONSchema := jsonSchemaValidatorMockData{
		times: 1,
		err:   nil,
	}

	invalidJSONSchema := jsonSchemaValidatorMockData{
		times: 1,
		err:   errors.New("bad-request").SetType(errors.ErrorTypeBadRequest),
	}

	nonExistingJSONSchema := jsonSchemaValidatorMockData{
		times: 1,
		err:   errors.New("some-error").SetType(errors.ErrorTypeInternal),
	}

	successfulAddDocument := domainServiceAddDocumentMockData{
		times: 1,
		err:   nil,
	}

	failedToAddDocument := domainServiceAddDocumentMockData{
		times: 1,
		err:   errors.New("some-error"),
	}

	tests := []struct {
		name                       string
		jsonSchemaValidatorMD      jsonSchemaValidatorMockData
		domainServiceAddDocumentMD domainServiceAddDocumentMockData
		body                       models.Document
		wantedStatusCode           int
		wantErr                    bool
	}{
		{
			name:                       "add document successfully expect status OK (200)",
			jsonSchemaValidatorMD:      validJSONSchema,
			domainServiceAddDocumentMD: successfulAddDocument,
			body:                       validDoc,
			wantedStatusCode:           http.StatusOK,
			wantErr:                    false,
		},
		{
			name:                  "invalid document reported expect status bad request (400)",
			jsonSchemaValidatorMD: invalidJSONSchema,
			body:                  invalidDoc,
			wantedStatusCode:      http.StatusBadRequest,
			wantErr:               true,
		},
		{
			name:                  "failed to validate reported document expect status internal server error (500)",
			jsonSchemaValidatorMD: nonExistingJSONSchema,
			wantedStatusCode:      http.StatusInternalServerError,
			wantErr:               true,
		},
		{
			name:                  "failed to validate reported document expect status internal server error (500)",
			jsonSchemaValidatorMD: nonExistingJSONSchema,
			wantedStatusCode:      http.StatusInternalServerError,
			wantErr:               true,
		},
		{
			name:                       "failed to add reported document to db expect status internal server error (500)",
			jsonSchemaValidatorMD:      validJSONSchema,
			domainServiceAddDocumentMD: failedToAddDocument,
			body:                       invalidDoc,
			wantedStatusCode:           http.StatusInternalServerError,
			wantErr:                    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			reportedDocumentInByte, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("Failed to marshal document (%+v) from request body. Error: %s", tt.body, err)
			}

			id := uuid.New().String()

			domainService := mocks.NewMockDomainService(c)
			domainService.EXPECT().AddDocument(gomock.Any(), tt.body).
				Times(tt.domainServiceAddDocumentMD.times).
				Return(id, tt.domainServiceAddDocumentMD.err)

			jsonSchemaValidator := mocks.NewMockJSONSchemaValidator(c)
			jsonSchemaValidator.EXPECT().ValidateSchemaFromBytes(postDocumentSchemaName, reportedDocumentInByte).
				Times(tt.jsonSchemaValidatorMD.times).
				Return(tt.jsonSchemaValidatorMD.err)

			s := &Adapter{
				domainSvc:  domainService,
				jsonSchema: jsonSchemaValidator,
			}

			r := chi.NewRouter()
			r.Route("/documents", func(r chi.Router) {
				r.Post("/", s.addDocument)
			})

			ts := httptest.NewServer(r)
			defer ts.Close()

			res, resBody := testRequest(t, ts, http.MethodPost, "/documents", bytes.NewReader(reportedDocumentInByte))
			statusCodeCheck(t, res, tt.wantedStatusCode)

			if tt.wantErr {
				return
			}

			if id != string(resBody) {
				t.Fatalf("addDocument() got = %s, want %v", string(resBody), id)
			}
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method string, path string, body io.Reader) (*http.Response, []byte) {
	url := ts.URL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("Failed to close response body. Error: %s", err)
		}
	}()

	return res, resBody
}

func statusCodeCheck(t *testing.T, r *http.Response, wantedStatusCode int) {
	if r.StatusCode != wantedStatusCode {
		t.Fatalf("handler return wrong status code: got %s want %s", http.StatusText(r.StatusCode), http.StatusText(wantedStatusCode))
	}
}
