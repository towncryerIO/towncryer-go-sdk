package apis

import (
	"errors"
	"fmt"
)

var (
	ErrRequestPayloadLimitReached = fmt.Errorf("request limit reached, max allowed is %d", MAX_REQUEST_PAYLOAD_SIZE)
	ErrPayloadCannotBeNil         = errors.New("payload cannot be nil")
)
