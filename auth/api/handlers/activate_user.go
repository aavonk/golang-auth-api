package handlers

import (
	"net/http"

	"github.com/todo-app/api/helpers"
)

/**
-- Workflow for activating a user: --

1. The user submits the plaintext activation token (which they just received in their email)
to the PUT /v1/users/activated endpoint.

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
func ActivateUser() http.HandlerFunc {
	return activateUser()
}

func activateUser() http.HandlerFunc {
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

	}
}
