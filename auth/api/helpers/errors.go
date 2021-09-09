package helpers

import (
	"net/http"

	"github.com/todo-app/pkg/logger"
)

var (
	internalSrvErrMsg  = "the server encountered a problem and could not process your request"
	notFoundMssg       = "the requested reouce could not be found"
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
func ServerErrReponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Error.Println(err)
	errResponse(w, r, http.StatusInternalServerError, internalSrvErrMsg)
}

func NotFoundErrResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Error.Printf("Not found error: %v", err)
	errResponse(w, r, http.StatusNotFound, notFoundMssg)
}

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
