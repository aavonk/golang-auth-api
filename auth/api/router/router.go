package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/todo-app/api/handlers"
)

func Get() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "OK", "time": time.Now().UTC().String()})
	})

	r.HandleFunc("/register", handlers.Register).Methods("POST")

	http.Handle("/", r)
	return r
}
