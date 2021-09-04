package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/todo-app/internal"
	"github.com/todo-app/internal/domain"
)

// func Register(userRespository domain.UserRepository) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		register(w, r)
// 	}
// }

// TODO: Extract http error handling into package
func register(repo domain.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// Search for existing user
		found := repo.GetByEmail(user.Email)

		// If found isn't an empty user domain struct, then a user was found
		// in the db by that email, return
		if !found.IsEmpty() {
			internal.ErrUnprocessableEntity(errors.New("account exists"), "account already exists").Send(w)
			return
		}

		// validate that the user obj fits the domain and passes the
		// requirements.
		user.Prepare()
		err = user.Validate()

		if err != nil {
			internal.ErrUnprocessableEntity(err, err.Error()).Send(w)
			return
		}

		// Reassign user variable to what is returned from creating a user in the DB
		user, err = repo.Create(&user)

		if err != nil {
			internal.ErrInternalServer(err, "Internal server error").Send(w)
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
		session.Options = &sessions.Options{
			// Path:     "/",
			MaxAge:   86400 * 7, // one week
			Secure:   true,
			HttpOnly: true,
		}

		if err != nil {
			internal.ErrInternalServer(err, "Internal server error").Send(w)
			return
		}

		session.Values["jwt"] = token
		err = session.Save(r, w)

		if err != nil {
			internal.ErrInternalServer(err, "Internal server error").Send(w)
			return
		}
		//token
		userResponse := user.ToHTTPResponse()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&userResponse)

	}

}

func Register(app *internal.App) http.HandlerFunc {
	return register(app.UserRepository)
}
