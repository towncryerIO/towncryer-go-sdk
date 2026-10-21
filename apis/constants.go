package apis

import "time"

const (
	API_HEADER_KEY        = "X-API-KEY"
	API_HEADER_SECRET_KEY = "X-API-SECRET"

	MAX_REQUEST_PAYLOAD_SIZE       = 32768      // 32kb
	MAX_BATCH_REQUEST_PAYLOAD_SIZE = 520 * 1024 // 520kb

	// DEFAULT_MAX_RETRIES and DEFAULT_RETRY_DELAY are used when a
	// config.Configuration does not specify MaxRetries / RetryDelay.
	DEFAULT_MAX_RETRIES = 3
	DEFAULT_RETRY_DELAY = 150 * time.Millisecond

	// DEFAULT_TIMEOUT is applied to the underlying http.Client when a
	// config.Configuration does not specify a Timeout. It guarantees requests
	// fail instead of hanging forever when the caller passes a context with
	// no deadline (e.g. context.Background()).
	DEFAULT_TIMEOUT = 10 * time.Second
)
