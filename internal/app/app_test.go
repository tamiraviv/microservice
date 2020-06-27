package app

import (
	"context"
	"reflect"
	"testing"
	"time"

	"microservice/internal/pkg/errors"
	"microservice/mocks"

	"github.com/golang/mock/gomock"
)

func TestNewApp(t *testing.T) {
	type getLogLevelMockData struct {
		times    int
		logLevel string
		err      error
	}

	getValidLogLevel := getLogLevelMockData{
		times:    1,
		logLevel: "info",
		err:      nil,
	}

	getInvalidLogLevel := getLogLevelMockData{
		times:    1,
		logLevel: "fake-level",
		err:      nil,
	}

	failedToGetLogLevel := getLogLevelMockData{
		times: 1,
		err:   errors.New("some-error"),
	}

	tests := []struct {
		name          string
		getLogLevelMD getLogLevelMockData
		wantErr       bool
	}{
		{
			name:          "valid creation expect no error",
			getLogLevelMD: getValidLogLevel,
			wantErr:       false,
		},
		{
			name:          "failed to get log level from configuration expect error",
			getLogLevelMD: failedToGetLogLevel,
			wantErr:       true,
		},
		{
			name:          "failed to get log level from configuration expect error",
			getLogLevelMD: failedToGetLogLevel,
			wantErr:       true,
		},
		{
			name:          "get invalid log level from configuration expect error",
			getLogLevelMD: getInvalidLogLevel,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			restServer := mocks.NewMockRestServer(c)
			conf := mocks.NewMockConfigurationService(c)
			conf.EXPECT().GetString(confKeyLogLevel).Times(tt.getLogLevelMD.times).Return(tt.getLogLevelMD.logLevel, tt.getLogLevelMD.err)

			got, err := NewApp(conf, restServer)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			want := &App{
				restServer: restServer,
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("NewApp() = %v, want %v", got, want)
			}
		})
	}
}

func TestApp_Start(t *testing.T) {
	type startRestServerMockData struct {
		err error
	}

	successfulStartRestServer := startRestServerMockData{
		err: nil,
	}

	failedToStartRestServer := startRestServerMockData{
		err: errors.New("some-error"),
	}

	tests := []struct {
		name              string
		startRestServerMD startRestServerMockData
		wantErr           bool
	}{
		{
			name:              "successful start application",
			startRestServerMD: successfulStartRestServer,
			wantErr:           false,
		},
		{
			name:              "server failed to start rest server expect error",
			startRestServerMD: failedToStartRestServer,
			wantErr:           true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			restServer := mocks.NewMockRestServer(c)
			restServer.EXPECT().Start().AnyTimes().Return(tt.startRestServerMD.err)

			a := &App{
				restServer: restServer,
			}

			appStartErrorChan := make(chan error, 1)
			go func() {
				appStartErrorChan <- a.Start()
			}()

			select {
			case <-time.After(3 * time.Second):
				if tt.wantErr {
					t.Fatalf("App.Start() succssed, wantErr %v", tt.wantErr)
				}
			case err := <-appStartErrorChan:
				if (err != nil) != tt.wantErr {
					t.Fatalf("App.Start() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestApp_Stop(t *testing.T) {
	type stopRestServerMockData struct {
		times int
		err   error
	}

	successfulStopRestServer := stopRestServerMockData{
		times: 1,
		err:   nil,
	}

	failedToStopRestServer := stopRestServerMockData{
		times: 1,
		err:   errors.New("some-error"),
	}

	tests := []struct {
		name             string
		stopRestServerMD stopRestServerMockData
		wantErr          bool
	}{
		{
			name:             "successful stop expect no error",
			stopRestServerMD: successfulStopRestServer,
			wantErr:          false,
		},
		{
			name:             "failed to stop rest server expect error",
			stopRestServerMD: failedToStopRestServer,
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			restServer := mocks.NewMockRestServer(c)
			restServer.EXPECT().Stop(gomock.Any()).Times(tt.stopRestServerMD.times).Return(tt.stopRestServerMD.err)

			a := &App{
				restServer: restServer,
			}
			if err := a.Stop(context.TODO()); (err != nil) != tt.wantErr {
				t.Errorf("App.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
