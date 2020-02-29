package viper

import (
	"time"

	"microservice/internal/pkg/errors"

	"github.com/spf13/viper"
)

type Service struct {
	v *viper.Viper
}

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

func (v *Service) GetString(key string) (string, error) {
	value := v.v.Get(key)
	if value == nil {
		return "", errors.Errorf("Could not find key (%s) in configuration", key)
	}

	switch s := value.(type) {
	case string:
		return s, nil
	default:
		return "", errors.Errorf("value %#v of type %T is not string", value, value)
	}
}

func (v *Service) GetInt(key string) (int, error) {
	value := v.v.Get(key)
	if value == nil {
		return 0, errors.Errorf("Could not find key (%s) in configuration", key)
	}

	switch i := value.(type) {
	case int:
		return i, nil
	default:
		return 0, errors.Errorf("value %#v of type %T is not int", i, i)
	}
}

func (v *Service) GetDuration(key string) (time.Duration, error) {
	value := v.v.Get(key)
	if value == nil {
		return 0, errors.Errorf("Could not find key (%s) in configuration", key)
	}

	switch d := value.(type) {
	case string:
		return time.ParseDuration(d)
	default:
		return 0, errors.Errorf("value %#v of type %T is not a string", d, d)
	}
}
