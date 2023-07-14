package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/michelheusschen/omada_exporter/pkg/config"
)

func createMockServer() (*http.Client, *httptest.Server) {
	mockServer := httptest.NewServer(nil)
	mockClient := &http.Client{Transport: mockServer.Client().Transport}

	return mockClient, mockServer
}

func TestGetDevicesEmpty(t *testing.T) {
	mockClient, mockServer := createMockServer()
	defer mockServer.Close()

	client := &Client{
		Config:     &config.Config{Host: mockServer.URL},
		omadaCID:   "example",
		SiteId:     "123",
		httpClient: mockClient,
	}

	mockServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/example/api/v2/loginStatus" {
			_, _ = w.Write([]byte(`{"errorCode":0,"msg":"Success.","result":{"login":true}}`))
			return
		}

		// Verify the request
		if r.Method != "GET" || r.URL.Path != "/example/api/v2/sites/123/devices" {
			t.Errorf("Unexpected request: %s %s", r.Method, r.URL.Path)
		}

		_, _ = w.Write([]byte(`{"errorCode":0,"msg":"Success.","result":[]}`))
	})

	if _, err := client.GetDevices(); err != nil {
		t.Fatalf("GetDevices returned an error: %v", err)
	}
}
