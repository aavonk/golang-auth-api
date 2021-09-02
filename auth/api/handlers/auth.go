package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/todo-app/internal"
	"github.com/todo-app/internal/domain"
)

// TODO: Extract http error handling into package
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		internal.ErrUnprocessableEntity(err, err.Error()).Send(w)
		return
	}

	var user domain.User

	err = json.Unmarshal(body, &user)

	if err != nil {
		internal.ErrUnprocessableEntity(err, err.Error()).Send(w)

		return
	}

	user.Prepare()
	err = user.Validate()

	if err != nil {
		internal.ErrUnprocessableEntity(err, err.Error()).Send(w)

		return
	}

	// Generate a JWT
	jwtClaims := internal.JWTClaims{
		UserId: user.ID,
		Email:  user.Email,
	}

	token, err := internal.NewToken(jwtClaims)

	if err != nil {
		internal.ErrUnprocessableEntity(err, err.Error()).Send(w)

		return
	}

	// Store the JWT on the session
	session, err := internal.SessionStore.Get(r, "user-session")

	if err != nil {
		internal.ErrInternalServer(err, err.Error()).Send(w)
		return
	}

	session.Values["jwt"] = token
	session.Save(r, w)
	//token
	userResponse := user.ToHTTPResponse()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&userResponse)

}
