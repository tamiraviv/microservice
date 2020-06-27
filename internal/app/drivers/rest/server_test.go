package rest

import (
	"context"
	"os"
	"path"
	"reflect"
	"runtime"
	"testing"
	"time"

	"microservice/internal/pkg/errors"
	"microservice/mocks"

	"github.com/golang/mock/gomock"
)

func TestNewServer(t *testing.T) {
	type getServerPortMockData struct {
		times int
		err   error
	}

	type getServerTimeoutMockData struct {
		times int
		err   error
	}

	type jsonSchemaMockData struct {
		times int
		err   error
	}

	successfulGetServerPort := getServerPortMockData{
		times: 1,
		err:   nil,
	}

	failedToGetServerPort := getServerPortMockData{
		times: 1,
		err:   errors.New("some-error"),
	}

	successfulGetServerTimeout := getServerTimeoutMockData{
		times: 1,
		err:   nil,
	}

	failedToGetServerTimeout := getServerTimeoutMockData{
		times: 1,
		err:   errors.New("some-error"),
	}

	successfulSetJSONSchema := jsonSchemaMockData{
		times: 1,
		err:   nil,
	}

	failedToSetJSONSchema := jsonSchemaMockData{
		times: 1,
		err:   errors.New("some-error"),
	}

	tests := []struct {
		name               string
		getServerPortMD    getServerPortMockData
		getServerTimeoutMD getServerTimeoutMockData
		setJSONSchema      jsonSchemaMockData
		wantErr            bool
	}{
		{
			name:               "valid creation expect no error",
			getServerPortMD:    successfulGetServerPort,
			getServerTimeoutMD: successfulGetServerTimeout,
			setJSONSchema:      successfulSetJSONSchema,
			wantErr:            false,
		},
		{
			name:            "failed to get server port from configuration expect error",
			getServerPortMD: failedToGetServerPort,
			wantErr:         true,
		},
		{
			name:               "failed to get server timeout from configuration expect error",
			getServerPortMD:    successfulGetServerPort,
			getServerTimeoutMD: failedToGetServerTimeout,
			wantErr:            true,
		},
		{
			name:               "failed to set JSON Schema expect error",
			getServerPortMD:    successfulGetServerPort,
			getServerTimeoutMD: successfulGetServerTimeout,
			setJSONSchema:      failedToSetJSONSchema,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			port := 8080
			timeout := 10 * time.Second

			domainService := mocks.NewMockDomainService(c)

			js := mocks.NewMockJSONSchemaValidator(c)
			js.EXPECT().SetSchemaFromBytes(postDocumentSchemaName, gomock.AssignableToTypeOf([]byte{})).Times(tt.setJSONSchema.times).Return(tt.setJSONSchema.err)

			conf := mocks.NewMockConfigurationService(c)
			conf.EXPECT().GetInt(serverPortKey).Times(tt.getServerPortMD.times).Return(port, tt.getServerPortMD.err)
			conf.EXPECT().GetDuration(serverTimeoutKey).Times(tt.getServerTimeoutMD.times).Return(timeout, tt.getServerTimeoutMD.err)

			_, filename, _, _ := runtime.Caller(0)
			dir := path.Join(path.Dir(filename), "../../../../")
			_ = os.Chdir(dir)

			got, err := NewServer(conf, domainService, js)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got.port, port) {
				t.Errorf("NewServer() port = %v, want %v", got, port)
			}

			if !reflect.DeepEqual(got.timeout, timeout) {
				t.Errorf("NewServer() timeout = %v, want %v", got, timeout)
			}

			if !reflect.DeepEqual(got.domainSvc, domainService) {
				t.Errorf("NewServer() domainSvc = %v, want %v", got, domainService)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	type listenAndServeMockData struct {
		times int
		err   error
	}

	successfulListenAndServe := listenAndServeMockData{
		times: 1,
		err:   nil,
	}

	failedToListenAndServe := listenAndServeMockData{
		times: 1,
		err:   errors.New("some-error"),
	}

	tests := []struct {
		name             string
		listenAndServeMD listenAndServeMockData
		wantErr          bool
	}{
		{
			name:             "successful start expect no error",
			listenAndServeMD: successfulListenAndServe,
			wantErr:          false,
		},
		{
			name:             "failed to listen and serve expect error",
			listenAndServeMD: failedToListenAndServe,
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			httpServer := mocks.NewMockHTTPServer(c)
			httpServer.EXPECT().ListenAndServe().Times(tt.listenAndServeMD.times).Return(tt.listenAndServeMD.err)

			s := &Adapter{
				server: httpServer,
			}
			if err := s.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Adapter.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Stop(t *testing.T) {
	type teardownDomainServiceMockData struct {
		times int
		err   error
	}

	successfulTeardownDomainService := teardownDomainServiceMockData{
		times: 1,
		err:   nil,
	}

	failedToTeardownDomainService := teardownDomainServiceMockData{
		times: 1,
		err:   errors.New("some-error"),
	}

	tests := []struct {
		name                    string
		teardownDomainServiceMD teardownDomainServiceMockData
		wantErr                 bool
	}{
		{
			name:                    "successful stop expect no error",
			teardownDomainServiceMD: successfulTeardownDomainService,
			wantErr:                 false,
		},
		{
			name:                    "failed to teardown domain service expect error",
			teardownDomainServiceMD: failedToTeardownDomainService,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			domainService := mocks.NewMockDomainService(c)
			domainService.EXPECT().Teardown(gomock.Any()).Times(tt.teardownDomainServiceMD.times).Return(tt.teardownDomainServiceMD.err)

			s := &Adapter{
				domainSvc: domainService,
			}

			if err := s.Stop(context.TODO()); (err != nil) != tt.wantErr {
				t.Errorf("Adapter.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
