package rest

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"microservice/internal/pkg/errors"
	"microservice/models"

	"github.com/prometheus/common/log"
)

const (
	serverBaseKey    = "server"
	serverPortKey    = serverBaseKey + ".port"
	serverTimeoutKey = serverBaseKey + ".timeout"

	apiFolder              = "api"
	postDocumentSchemaName = "PostDocument"
	postDocumentSchemaFile = apiFolder + "/" + "postDocumentSchema.json"
)

// Configuration expose an interface of configuration related actions
type Configuration interface {
	GetString(key string) (string, error)
	GetInt(key string) (int, error)
	GetDuration(key string) (time.Duration, error)
}

// DomainSvc exposes an interface of document related actions
type DomainSvc interface {
	GetDocument(ctx context.Context, id string) (models.Document, error)
	AddDocument(ctx context.Context, doc models.Document) (string, error)
	Teardown(ctx context.Context) error
}

// Server http server interface
type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// JSONSchemaValidator handle input json validation
type JSONSchemaValidator interface {
	SetSchemaFromBytes(name string, inputJSON []byte) error
	ValidateSchemaFromBytes(name string, inputJSON []byte) error
}

// Adapter defines the server struct
type Adapter struct {
	port       int
	timeout    time.Duration
	server     Server
	domainSvc  DomainSvc
	jsonSchema JSONSchemaValidator
}

// NewServer returns a new instance of the Adapter struct
func NewServer(conf Configuration, dsv DomainSvc, js JSONSchemaValidator) (*Adapter, error) {
	port, err := conf.GetInt(serverPortKey)
	if err != nil {
		return nil, err
	}

	timeout, err := conf.GetDuration(serverTimeoutKey)
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr: ":" + strconv.Itoa(port),
	}

	if err := initJSONSchemaValidator(js); err != nil {
		return nil, errors.Wrap(err, "Failed to init json schema validator")
	}

	a := &Adapter{
		port:       port,
		timeout:    timeout,
		server:     server,
		domainSvc:  dsv,
		jsonSchema: js,
	}

	server.Handler = a.newRouter(timeout)

	return a, nil
}

// Start REST web server
func (s *Adapter) Start() error {
	log.Infof("Starting rest server. Listening on port (%d)", s.port)
	if err := s.server.ListenAndServe(); err != nil {
		return errors.Wrapf(err, "failed to listen on port (%d)", s.port)
	}

	return nil
}

// Stop REST web server
func (s *Adapter) Stop(ctx context.Context) error {
	if err := s.domainSvc.Teardown(ctx); err != nil {
		return errors.Wrap(err, "Fail to tear down domain service")
	}

	return nil
}

func initJSONSchemaValidator(js JSONSchemaValidator) error {
	if err := calculateSchemaPath(apiFolder); err != nil {
		return errors.Wrap(err, "Failed to update reference in json schema files")
	}

	postDocumentSchema, err := ioutil.ReadFile(postDocumentSchemaFile)
	if err != nil {
		return errors.Wrap(err, "Failed to read post document schema")
	}

	if err := js.SetSchemaFromBytes(postDocumentSchemaName, postDocumentSchema); err != nil {
		return errors.Wrap(err, "Failed to set post document schema")
	}

	return nil
}

func calculateSchemaPath(subfix string) error {
	p, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Could not get absolute path of json schemas")
	}
	p = strings.ReplaceAll(path.Join(p, subfix), "\\", "/")

	files, err := filepath.Glob(p + "/*.json")
	if err != nil {
		return errors.Wrap(err, "Could not load all json schema files")
	}

	ptr := "(\\\"\\$ref\\\"\\:\\s*\"file:\\/\\/\\/)(((.*)(\\/.*#.*\"))|(([^#\"]*)(\\/[^#\"\\/]*\")))"
	rgx := regexp.MustCompile(ptr)

	for _, file := range files {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.Wrapf(err, "Could not read json schema file: %s", file)
		}

		fileAsString := string(f)
		str := rgx.ReplaceAllString(fileAsString, "${1}"+p+"$5$8")
		if err := ioutil.WriteFile(file, []byte(str), 0); err != nil {
			return errors.Wrapf(err, "Could not write to json schema file (%s) the content (%s)", file, str)
		}
	}

	return nil
}
