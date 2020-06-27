package jsonschema

import (
	"microservice/internal/pkg/errors"

	"github.com/xeipuuv/gojsonschema"
)

// Service contains map of json schemas
type Service struct {
	schemas map[string]*gojsonschema.Schema
}

// NewJSONSchemaService returns a new instance of the Service struct
func NewJSONSchemaService() *Service {
	return &Service{
		schemas: make(map[string]*gojsonschema.Schema),
	}
}

// SetSchemaFromString add new json schema from string
func (s *Service) SetSchemaFromString(name string, inputJSON string) error {
	schemaLoader := gojsonschema.NewStringLoader(inputJSON)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return errors.Errorf("Failed to read json schema: %s", err).SetType(errors.ErrorTypeBadRequest)
	}

	s.schemas[name] = schema
	return nil
}

// SetSchemaFromBytes add new json schema from bytes
func (s *Service) SetSchemaFromBytes(name string, inputJSON []byte) error {
	schemaLoader := gojsonschema.NewBytesLoader(inputJSON)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return errors.Errorf("Failed to read json schema: %s", err).SetType(errors.ErrorTypeBadRequest)
	}

	s.schemas[name] = schema
	return nil
}

// ValidateSchemaFromString validate string json input with schema
func (s *Service) ValidateSchemaFromString(name string, inputJSON string) error {
	baseSchema, ok := s.schemas[name]
	if !ok {
		return errors.Errorf("No schema found for: %s", name).SetType(errors.ErrorTypeInternal)
	}

	v, err := baseSchema.Validate(gojsonschema.NewStringLoader(inputJSON))
	if err != nil {
		return errors.Wrap(err, "Failed to read json input").SetType(errors.ErrorTypeBadRequest)
	}

	if !v.Valid() {
		return errors.Errorf("Json is not valid according to the JsonSchema. Errors: %s", v.Errors()).SetType(errors.ErrorTypeBadRequest)
	}

	return nil
}

// ValidateSchemaFromBytes validate bytes json input with schema
func (s *Service) ValidateSchemaFromBytes(name string, inputJSON []byte) error {
	baseSchema, ok := s.schemas[name]
	if !ok {
		return errors.Errorf("No schema found for: %s", name).SetType(errors.ErrorTypeInternal)
	}

	v, err := baseSchema.Validate(gojsonschema.NewBytesLoader(inputJSON))
	if err != nil {
		return errors.Wrap(err, "Failed to read json input").SetType(errors.ErrorTypeBadRequest)
	}

	if !v.Valid() {
		return errors.Errorf("Json is not valid according to the JsonSchema. Errors: %s", v.Errors()).SetType(errors.ErrorTypeBadRequest)
	}

	return nil
}
