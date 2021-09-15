package handlers

import (
	"errors"
	"net/http"

	"github.com/todo-app/api/helpers"
	"github.com/todo-app/internal/application"
	"github.com/todo-app/internal/domain"
	"github.com/todo-app/internal/repositories"
	"github.com/todo-app/internal/validator"
)

/**
-- Workflow for activating a user: --

1. The user submits the plaintext activation token (which they just received in their email)
to the PUT /v1/users/activate endpoint.

2. We validate the plaintext token to check that it matches the expected format, sending the
client an error message if necessary.

3. We then call the UserRepository.GetForToken() method to retrieve the details of the user associated
with the provided token. If there is no matching token found, or it has expired, we send the client
an error message.

4. We activate the associated user by setting activated = true on the user record and update it in our database.

5. We delete all activation tokens for the user from the tokens table. We can do this using the
TokenRepository.DeleteAllForUser() method that we made earlier.

6. We send the updated user details in a JSON response.
*/
func ActivateUser(app *application.App) http.HandlerFunc {

	return activateUser(app.UserRepository, app.TokenRepository)
}

func activateUser(userRepo repositories.UserRepositoryInterface, tokenRepo repositories.TokenRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse the plaintext activation token from the request
		var input struct {
			TokenPlainText string `json:"token"`
		}

		err := helpers.ReadJSON(w, r, &input)
		if err != nil {
			helpers.BadRequestErrResponse(w, r, err)
			return
		}

		v := validator.New()

		if domain.ValidateTokenPlainText(v, input.TokenPlainText); !v.Valid() {
			helpers.FailedValidationResponse(w, r, v.Errors)
			return
		}

		user, err := userRepo.GetForToken(domain.TokenScopeActivation, input.TokenPlainText)
		if err != nil {
			switch {
			case errors.Is(err, repositories.ErrRecordNotFound):
				v.AddError("token", "invalid or expired activation token")
				helpers.FailedValidationResponse(w, r, v.Errors)
			default:
				helpers.ServerErrReponse(w, r, err)

			}
			return
		}

		// Update the user's activation status
		user.Activated = true

		// Save the updated record in the database, checking for any edit conflicts
		err = userRepo.Update(user)
		if err != nil {
			switch {
			case errors.Is(err, repositories.ErrEditConflict):
				helpers.UnprocessableErrResponse(w, r, err)
			default:
				helpers.ServerErrReponse(w, r, err)
			}
			return
		}

		// If everything went successfully, delete all activation tokens
		// associated with the user.
		err = tokenRepo.DeleteAllForUser(domain.TokenScopeActivation, user.ID.String())
		if err != nil {
			helpers.ServerErrReponse(w, r, err)
			return
		}

		err = helpers.SendJSON(w, http.StatusOK, user.ToHTTPResponse(), nil)
		if err != nil {
			helpers.ServerErrReponse(w, r, err)
		}
	}
}
