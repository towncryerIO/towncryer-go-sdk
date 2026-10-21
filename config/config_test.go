package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/towncryerIO/towncryer-go-sdk/config"
)

func TestConfiguration_WithEnvDefaults(t *testing.T) {
	t.Run("fills blank fields from environment variables", func(t *testing.T) {
		t.Setenv("TOWNCRYER_API_KEY", "env-api-key")
		t.Setenv("TOWNCRYER_API_SECRET", "env-api-secret")
		t.Setenv("TOWNCRYER_BASE_URL", "https://env.example.com")

		cfg := config.Configuration{}.WithEnvDefaults()

		assert.Equal(t, "env-api-key", cfg.ApiKey)
		assert.Equal(t, "env-api-secret", cfg.ApiSecret)
		assert.Equal(t, "https://env.example.com", cfg.BaseUrl)
	})

	t.Run("does not override explicitly set fields", func(t *testing.T) {
		t.Setenv("TOWNCRYER_API_KEY", "env-api-key")
		t.Setenv("TOWNCRYER_API_SECRET", "env-api-secret")
		t.Setenv("TOWNCRYER_BASE_URL", "https://env.example.com")

		cfg := config.Configuration{
			ApiKey:    "explicit-api-key",
			ApiSecret: "explicit-api-secret",
			BaseUrl:   "https://explicit.example.com",
		}.WithEnvDefaults()

		assert.Equal(t, "explicit-api-key", cfg.ApiKey)
		assert.Equal(t, "explicit-api-secret", cfg.ApiSecret)
		assert.Equal(t, "https://explicit.example.com", cfg.BaseUrl)
	})
}
