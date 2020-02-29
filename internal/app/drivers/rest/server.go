package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"microservice/internal/pkg/errors"
	"microservice/models"

	"github.com/prometheus/common/log"
)

const (
	serverBaseKey    = "server"
	serverPortKey    = serverBaseKey + ".port"
	serverTimeoutKey = serverBaseKey + ".timeout"
)

type Configuration interface {
	GetString(key string) (string, error)
	GetInt(key string) (int, error)
	GetDuration(key string) (time.Duration, error)
}

type DomainSvc interface {
	GetDocument(ctx context.Context, id string) (models.Document, error)
	AddDocument(ctx context.Context, doc models.Document) (string, error)
	TearDown(ctx context.Context) error
}

type Server struct {
	port      int
	timeout   time.Duration
	domainSvc DomainSvc
}

func NewServer(conf Configuration, dsv DomainSvc) (*Server, error) {
	port, err := conf.GetInt(serverPortKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get port from configuration")
	}

	timeout, err := conf.GetDuration(serverTimeoutKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get timeout from configuration")
	}

	return &Server{
		port:      port,
		timeout:   timeout,
		domainSvc: dsv,
	}, nil
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	r := s.newRouter(s.timeout)
	log.Infof("Starting rest server. Listening on port (%d)", s.port)
	if err := http.ListenAndServe(addr, r); err != nil {
		return errors.Wrapf(err, "failed to listen on port (%d)", s.port)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.domainSvc.TearDown(ctx); err != nil {
		return errors.Wrap(err, "Fail to tear down domain service")
	}

	return nil
}
