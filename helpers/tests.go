package helpers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volvlabs/towncryer-sdk/config"
)

type MockResponse struct {
	StatusCode int
	Payload    []byte
}

func CreateHttpServer(t *testing.T, endpoint string, statusCode int, payload []byte) (*config.Configuration, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, endpoint, r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		w.WriteHeader(statusCode)
		//nolint:errcheck
		w.Write(payload)
	}))

	config := &config.Configuration{
		BaseUrl:   server.URL,
		ApiKey:    "dummy_api_key",
		ApiSecret: "dummy_api_secret",
	}

	return config, server
}
