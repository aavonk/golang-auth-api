package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/todo-app/api/helpers"
	"github.com/todo-app/pkg/logger"

	"github.com/todo-app/internal/application"
	"github.com/todo-app/internal/domain"
	"github.com/todo-app/internal/mailer"
	"github.com/todo-app/internal/repositories"
	"github.com/todo-app/internal/services"
	"github.com/todo-app/internal/validator"
)

func register(service services.IdentityServiceInterface, tokenRepo repositories.TokenRepositoryInterface, mailer mailer.Mailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user domain.User

		err := helpers.ReadJSON(w, r, &user)

		if err != nil {
			helpers.BadRequestErrResponseWithMsg(w, r, err)
			return
		}

		user.Prepare()
		v := validator.New()

		v.Check(user.FirstName != "", "firstName", "first name is required")
		v.Check(user.LastName != "", "lastName", "last name is required")
		v.Check(v.Matches(user.Email, validator.EmailRX), "email", "invalid email")
		v.Check(len([]rune(user.Password)) >= 6, "password", "password must be atleast 6 characters")

		if !v.Valid() {
			helpers.FailedValidationResponse(w, r, v.Errors)
			return
		}

		createdUser, err := service.HandleRegister(&user)

		if err != nil {
			switch {
			case errors.Is(err, repositories.ErrDuplicateEmail):
				helpers.BadRequestErrResponseWithMsg(w, r, err)
			default:
				helpers.ServerErrReponse(w, r, err)
			}
			return
		}

		// After the user record has been created in the database, generate a new activation
		// token for the user.
		token, err := tokenRepo.New(createdUser.ID.String(), 3*24*time.Hour, domain.TokenScopeActivation)
		if err != nil {
			helpers.ServerErrReponse(w, r, err)
			return
		}

		userResponse := createdUser.ToHTTPResponse()

		// Send Welcome email in a goroutine so it gets processed in the background
		go func() {
			// Run a deferred function which uses recover() to catch any panic, and log an
			// error message instead of terminating the application.
			defer func() {
				if err := recover(); err != nil {
					logger.Error.Println(fmt.Errorf("%s", err))
				}
			}()

			data := map[string]interface{}{
				"activationToken": token.Plaintext,
				"user":            createdUser,
			}

			err = mailer.Send(createdUser.Email, "user_welcome.tmpl", data)
			if err != nil {
				// We just want to log it instead of send an error response
				// to the client
				logger.Error.Println(err)
			}
		}()
		// Set the header as 202 Accepted to indicate that the request
		// has been accepted for processing but may not be complete because
		// the email could still be sending.
		err = helpers.SendJSON(w, http.StatusAccepted, userResponse, nil)
		if err != nil {
			helpers.ServerErrReponse(w, r, err)
		}

	}

}

func Register(app *application.App) http.HandlerFunc {
	return register(app.IdentityService, app.TokenRepository, app.Mailer)
}
