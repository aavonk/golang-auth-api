package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/todo-app/internal/domain"
)

// TODO: Extract http error handling into package
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var user domain.User

	err = json.Unmarshal(body, &user)

	user.Prepare()

	fmt.Printf("User struct: %+v", user)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	err = user.Validate()

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	userResponse := user.ToHTTPResponse()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&userResponse)

}
