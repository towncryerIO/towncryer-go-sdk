package towncryer

import "github.com/towncryerIO/towncryer-go-sdk/apis"

// APIError represents a non-2xx response from the Towncryer API. See apis.APIError
// for details. It is aliased here so callers can write:
//
//	var apiErr *towncryer.APIError
//	if errors.As(err, &apiErr) {
//	    switch apiErr.StatusCode {
//	    case 401: ...
//	    case 400: ...
//	    }
//	}
type APIError = apis.APIError
