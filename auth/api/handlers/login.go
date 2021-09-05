package handlers

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	"net/http"

// 	"github.com/todo-app/internal"
// 	"github.com/todo-app/internal/domain"
// )

// func Login(app *internal.App) http.HandlerFunc {
// 	return login(app.UserRepository)
// }

// func login(repo domain.UserRepository) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		body, err := ioutil.ReadAll(r.Body)

// 		if err != nil {
// 			internal.ErrUnprocessableEntity(err, "cannot parse body").Send(w)
// 			return
// 		}

// 		var loginReq loginRequest

// 		err = json.Unmarshal(body, &loginReq)

// 		if err != nil {
// 			internal.ErrUnprocessableEntity(err, "cannot parse body").Send(w)
// 			return
// 		}

// 		// Pass it off to authservice.HandleLogin()
// 	}
// }
