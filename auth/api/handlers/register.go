package handlers

import (
	"fmt"
	"net/http"

	"github.com/todo-app/api/helpers"
	"github.com/todo-app/pkg/logger"

	"github.com/todo-app/internal/application"
	"github.com/todo-app/internal/domain"
	"github.com/todo-app/internal/identity"
	"github.com/todo-app/internal/mailer"
	"github.com/todo-app/internal/services"
)

func register(service services.IdentityServiceInterface, mailer mailer.Mailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user domain.User

		err := helpers.ReadJSON(w, r, &user)

		if err != nil {
			helpers.BadRequestErrResponseWithMsg(w, r, err)
			return
		}

		createdUser, err := service.HandleRegister(&user)

		if err != nil {
			helpers.BadRequestErrResponseWithMsg(w, r, err)
			return
		}

		userResponse := createdUser.ToHTTPResponse()

		err = identity.SetCookie(w, &user)
		if err != nil {
			helpers.ServerErrReponse(w, r, err)

			return
		}
		// Send Welcome email in a goroutine so it gets processed in the background
		go func() {
			// Run a deferred function which uses recover() to catch any panic, and log an
			// error message instead of terminating the application.
			defer func() {
				if err := recover(); err != nil {
					logger.Error.Println(fmt.Errorf("%s", err))
				}
			}()
			err = mailer.Send(user.Email, "user_welcome.tmpl", user)
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
	return register(app.IdentityService, app.Mailer)
}
