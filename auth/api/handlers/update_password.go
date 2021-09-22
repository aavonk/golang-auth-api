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

func UpdateUserPasswordHandler(app *application.App) http.HandlerFunc {
	return updateUserPasswordHandler(app.UserRepository, app.TokenRepository)
}

// Verify the password reset token and set a new password for the user.
func updateUserPasswordHandler(userRepo repositories.UserRepositoryInterface, tokenRepo repositories.TokenRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Password       string `json:"password"`
			TokenPlaintext string `json:"token"`
		}

		err := helpers.ReadJSON(w, r, &input)
		if err != nil {
			helpers.BadRequestErrResponse(w, r, err)
			return
		}
		v := validator.New()

		domain.ValidateTokenPlainText(v, input.TokenPlaintext)
		v.Check(len([]rune(input.Password)) >= 6, "password", "password must be atleast 6 characters")

		if !v.Valid() {
			helpers.FailedValidationResponse(w, r, v.Errors)
			return
		}

		// Retrieve the details of the user associated with the password reset token,
		// returning an error message if no matching record was found.
		user, err := userRepo.GetForToken(domain.TokenScopePasswordReset, input.TokenPlaintext)
		if err != nil {
			switch {
			case errors.Is(err, repositories.ErrRecordNotFound):
				v.AddError("token", "invalid or expired token")
				helpers.FailedValidationResponse(w, r, v.Errors)
			default:
				helpers.ServerErrReponse(w, r, err)
			}
			return
		}

		// Set the new password for the user and hash it
		user.Password = input.Password
		user.HashPassword()
		// Save the updated user record in our database, checking for any edit conflicts as
		// normal.
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

		// If everything was successful, delete all password reset tokens for the user
		err = tokenRepo.DeleteAllForUser(domain.TokenScopePasswordReset, user.ID.String())
		if err != nil {
			helpers.ServerErrReponse(w, r, err)
			return
		}

		response := map[string]interface{}{
			"success": true,
			"message": "password successfully reset",
		}
		err = helpers.SendJSON(w, http.StatusOK, response, nil)
		if err != nil {
			helpers.ServerErrReponse(w, r, err)
		}
	}
}
