package errors

import "errors"

// ErrServerUnavailable is an error variable that represents a situation where
// the server is not currently available to handle a request.
var ErrServerUnavailable = errors.New("server unavailable")
