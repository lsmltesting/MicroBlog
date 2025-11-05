package http

import (
	"net/http"
)

// HTTPHandler interface is used to create handlers for each model
type HTTPHandler interface {
	sendError(w http.ResponseWriter, message string, statusCode int)
}
