package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"microservice/internal/pkg/errors"
	"microservice/models"

	"github.com/go-chi/chi"
)

const (
	URLParamID = "id"
)

func (s *Server) getDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, URLParamID)
	doc, err := s.domainSvc.GetDocument(ctx, id)
	if err != nil {
		if errors.IsType(err, errors.ErrorTypeNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		} else if errors.IsType(err, errors.ErrorTypeBadRequest) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(doc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (s *Server) addDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var doc models.Document
	if err := json.Unmarshal(body, &doc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := s.domainSvc.AddDocument(ctx, doc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(id))
	w.WriteHeader(http.StatusOK)
}
