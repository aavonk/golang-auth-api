package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/todo-app/api/helpers"
	"github.com/todo-app/internal/application"
	"github.com/todo-app/internal/domain"
	"github.com/todo-app/internal/mailer"
	"github.com/todo-app/internal/repositories"
	"github.com/todo-app/internal/validator"
	"github.com/todo-app/pkg/logger"
)

/** Workflow for password reset:

1. A client sends a request to the POST /v1/tokens/password-reset endpoint containing the email address
of the user whose password they want to reset.

2. If a user with that email address exists in the users table, and the user has already confirmed their
   email address by activating, then generate a cryptographically-secure high-entropy random token.

3. Store a hash of this token in the tokens table, alongside the user ID and a short (30-60 minute) expiry time for the token.

4. Send the original (unhashed) token to the user in an email.

5. If the owner of the email address didn’t request a password reset token, they can ignore the email.

6. Otherwise, they can submit the token to the PUT /v1/users/password endpoint along with their new password.
   If the hash of the token exists in the tokens table and hasn’t expired, then generate a bcrypt hash of the
   new password and update the user’s record.

7 Delete all existing password reset tokens for the user.

*/

func PasswordReset(app *application.App) http.HandlerFunc {
	return passwordReset(app.UserRepository, app.TokenRepository, app.Mailer)
}

func passwordReset(userRepo repositories.UserRepositoryInterface, tokenRepo repositories.TokenRepositoryInterface, mailer mailer.Mailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Email string `json:"email"`
		}

		err := helpers.ReadJSON(w, r, &input)
		if err != nil {
			helpers.BadRequestErrResponse(w, r, err)
			return
		}

		v := validator.New()
		v.Check(v.Matches(input.Email, validator.EmailRX), "email", "invalid email")
		if !v.Valid() {
			helpers.FailedValidationResponse(w, r, v.Errors)
			return
		}

		// Try to retrieve the corresponding user record for the email address. If it can't
		// be found, return an error message to the client.
		user, err := userRepo.GetByEmail(input.Email)
		if err != nil {
			switch {
			case errors.Is(err, repositories.ErrRecordNotFound):
				v.AddError("email", "no matching email found")
				helpers.FailedValidationResponse(w, r, v.Errors)
			default:
				helpers.ServerErrReponse(w, r, err)
			}
			return
		}

		// Return an error message if the user hasn't activated their account.
		if !user.Activated {
			v.AddError("email", "an account must be activated")
			helpers.FailedValidationResponse(w, r, v.Errors)
			return
		}

		// Otherwise, create a new password reset token with a 45-minute expiry time.
		token, err := tokenRepo.New(user.ID.String(), 45*time.Minute, domain.TokenScopePasswordReset)
		if err != nil {
			helpers.ServerErrReponse(w, r, err)
			return
		}

		// Email the user with their password reset token. Send the email in a background
		// go routine
		go func() {
			// Handle any errors from this goroutine as it wont be caught from the
			// panic recovery middleware
			defer func() {
				if err := recover(); err != nil {
					logger.Error.Println(fmt.Errorf("%s", err))
				}
			}()

			data := map[string]interface{}{
				"passwordResetToken": token.Plaintext,
			}

			mailer.Send(user.Email, "password_reset.tmpl", data)
			if err != nil {
				// We just want to log it instead of send an error response
				// to the client as they might have already gotten the response below.
				logger.Error.Println(err)
			}
		}()

		response := map[string]interface{}{
			"success": true,
			"message": "an email will be sent to you containing password reset instructions",
		}

		err = helpers.SendJSON(w, http.StatusAccepted, response, nil)
		if err != nil {
			helpers.ServerErrReponse(w, r, err)
		}

	}
}
