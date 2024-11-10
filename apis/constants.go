package apis

import "time"

const (
	API_HEADER_KEY        = "X-API-KEY"
	API_HEADER_SECRET_KEY = "X-API-SECRET"

	MAX_REQUEST_PAYLOAD_SIZE       = 32768      // 32kb
	MAX_BATCH_REQUEST_PAYLOAD_SIZE = 520 * 1024 // 520kb

	MAX_NO_OF_RETRIES = 3
	RETRY_DELAY       = 150 * time.Millisecond
)
