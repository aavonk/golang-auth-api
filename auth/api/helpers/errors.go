package helpers

import (
	"fmt"
	"net/http"

	"github.com/todo-app/pkg/logger"
)

var (
	internalSrvErrMsg  = "the server encountered a problem and could not process your request"
	notFoundMssg       = "the requested resource could not be found"
	unproccessagbleMsg = "the given data was not processable"
	unauthorizedMsg    = "unauthorized"
)

func errResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	//Write the response using the writeJSON helper. If this returns an error, log it, and fall back to
	// sending the client an empty response with a 500 Internal Server Error status code
	err := SendJSON(w, status, env, nil)
	if err != nil {
		logger.Error.Println(err)
		w.WriteHeader(status)
	}
}

// ServerErrResponse writes a Internal Service Error 500 Status Code to the response writer
// and writes back an error message of "The server encountered a problem and could not process your request"
func ServerErrReponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Error.Println(err)
	errResponse(w, r, http.StatusInternalServerError, internalSrvErrMsg)
}

// NotFoundErrResponse writes a Status Not Found 404 Status code to the response writer
// and writes back an error message of "the requested resource could not be found"
func NotFoundErrResponse(w http.ResponseWriter, r *http.Request) {
	logger.Error.Println("Not found error")
	errResponse(w, r, http.StatusNotFound, notFoundMssg)
}

// UnprocessableErrResponse write a Status Code of 422 - StatusUnprocessableEntity with the message
// "the given data was not processable"
func UnprocessableErrResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Error.Println(err)
	errResponse(w, r, http.StatusUnprocessableEntity, unproccessagbleMsg)
}

func UnauthorizedErrResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Error.Printf("UNAUTHORIZED - %v", err)
	errResponse(w, r, http.StatusUnauthorized, unauthorizedMsg)
}

func InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Error.Println(err)
	errResponse(w, r, http.StatusUnauthorized, "invalid credentials")
}

// MethodNotAllowedResponse will be used to send a 405 Method Not Allowed
// status code and JSON response to the client.
func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	errResponse(w, r, http.StatusMethodNotAllowed, message)
}
