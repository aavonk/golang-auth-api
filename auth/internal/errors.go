package internal

import (
	"encoding/json"
	"net/http"

	"github.com/todo-app/internal/logger"
)

type ErrResponse struct {
	Err            error  `json:"-"`               // low level runtime error
	HTTPStatusCode int    `json:"-"`               // http response status code
	Message        string `json:"message"`         // user-level status message
	ErrorText      string `json:"error,omitempty"` // application-level error message for debugging
}

func ErrUnprocessableEntity(err error, message string) *ErrResponse {
	logger.Error.Printf("Unprocessable Entity -- Received error: %+v . User facing error message: %v", err, message)
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		Message:        message,
		ErrorText:      err.Error(),
	}
}

func ErrInternalServer(err error, message string) *ErrResponse {
	logger.Error.Printf("Internal Server Error -- Received error: %v . User facing error message: %v", err, message)

	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		Message:        message,
		ErrorText:      err.Error(),
	}
}

func (e *ErrResponse) Send(w http.ResponseWriter) {
	w.WriteHeader(e.HTTPStatusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": e.Message})
}
