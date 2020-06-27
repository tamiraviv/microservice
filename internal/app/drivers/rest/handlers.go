package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"microservice/internal/pkg/errors"
	"microservice/models"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

const (
	urlParamID = "id"
)

func (s *Adapter) getDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, urlParamID)
	doc, err := s.domainSvc.GetDocument(ctx, id)
	if err != nil {
		if errors.IsType(err, errors.ErrorTypeNotFound) {
			log.Debugf("Could not found document with id (%s)", id)
			returnHTTPError(w, http.StatusNotFound, err.Error())
			return
		} else if errors.IsType(err, errors.ErrorTypeBadRequest) {
			log.Debugf("Failed to get document with id (%s) from domain. Error: %s", id, err)
			returnHTTPError(w, http.StatusBadRequest, "Invalid request id")
			return
		}

		log.Errorf("Failed to get document with id (%s) from domain. Error: %s", id, err)
		returnHTTPError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	b, err := json.Marshal(doc)
	if err != nil {
		log.Errorf("Failed to marshal document (%+v). Error: %s", doc, err)
		returnHTTPError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	httpReturn(w, http.StatusOK, b)
}

func (s *Adapter) addDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to read request body. Error: %s", err)
		returnHTTPError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if err := s.jsonSchema.ValidateSchemaFromBytes(postDocumentSchemaName, body); err != nil {
		if errors.IsType(err, errors.ErrorTypeBadRequest) {
			log.Debugf("Invalid schema: %s", err)
			returnHTTPError(w, http.StatusBadRequest, err.Error())
		} else {
			log.Errorf("Failed to validate request body: %s", err)
			returnHTTPError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		return
	}

	var doc models.Document
	if err := json.Unmarshal(body, &doc); err != nil {
		log.Debugf("Failed to unmarshal document. Error: %s", err)
		returnHTTPError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	id, err := s.domainSvc.AddDocument(ctx, doc)
	if err != nil {
		log.Errorf("Failed to validate request body: %s", err)
		returnHTTPError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	httpReturn(w, http.StatusOK, []byte(id))
}

func httpReturn(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(body); err != nil {
		log.Error("Failed to write HTTP response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}

func returnHTTPError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(msg)); err != nil {
		log.Error("Failed to write HTTP response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
