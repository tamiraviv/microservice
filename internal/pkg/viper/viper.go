package viper

import (
	"time"

	"microservice/internal/pkg/errors"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// Service implements the configuration service
type Service struct {
	v *viper.Viper
}

// NewConfiguration returns a new instance of the Service struct
func NewConfiguration() (*Service, error) {
	v := viper.New()
	v.SetConfigFile("./conf/bootstrapConfiguration.yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Failed to load configuration")
	}

	return &Service{
		v: v,
	}, nil
}

// IsSet checks if the requested key exists
func (v *Service) IsSet(key string) bool {
	return v.v.IsSet(key)
}

// Set places the key and value in configuration service
func (v *Service) Set(key string, value interface{}) {
	v.v.Set(key, value)
}

// Get returns the value for the requested key as interface{}
func (v *Service) Get(key string) interface{} {
	return v.v.Get(key)
}

// GetString returns the value for the requested key as string
func (v *Service) GetString(key string) (string, error) {
	if !v.v.IsSet(key) {
		return "", keyNotFoundError(key)
	}

	value, err := cast.ToStringE(v.v.Get(key))
	if err != nil {
		return "", invalidDataTypeError(key)
	}

	return value, nil
}

// GetInt returns the value for the requested key as int
func (v *Service) GetInt(key string) (int, error) {
	if !v.v.IsSet(key) {
		return 0, keyNotFoundError(key)
	}

	value, err := cast.ToIntE(v.v.Get(key))
	if err != nil {
		return 0, invalidDataTypeError(key)
	}

	return value, nil
}

// GetDuration returns the value for the requested key as duration
func (v *Service) GetDuration(key string) (time.Duration, error) {
	if !v.v.IsSet(key) {
		return 0, keyNotFoundError(key)
	}

	value, err := time.ParseDuration(v.v.GetString(key))
	if err != nil {
		return 0, invalidDataTypeError(key)
	}

	return value, nil
}

// GetBool returns the value for the requested key as bool
func (v *Service) GetBool(key string) (bool, error) {
	if !v.v.IsSet(key) {
		return false, keyNotFoundError(key)
	}

	value, err := cast.ToBoolE(v.v.Get(key))
	if err != nil {
		return false, invalidDataTypeError(key)
	}

	return value, nil
}

func keyNotFoundError(key string) error {
	return errors.Errorf("Failed to get key (%s) from configuration", key).SetType(errors.ErrorTypeNotFound)
}

func invalidDataTypeError(key string) error {
	return errors.Errorf("Invalid data type for key (%s) from configuration", key).SetType(errors.ErrorTypeInternal)
}
