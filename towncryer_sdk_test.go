package towncryer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volvlabs/towncryer-sdk"
	"github.com/volvlabs/towncryer-sdk/config"
)

func TestTowncryerSdk(t *testing.T) {
	// Arrange:
	cfg := &config.Configuration{}
	towncryerSdk := towncryer.New(cfg)

	// Act-Assert:

	assert.NotNil(t, towncryerSdk.Customers)
	assert.NotNil(t, towncryerSdk.Events)
	assert.NotNil(t, towncryerSdk.Messages)
}
