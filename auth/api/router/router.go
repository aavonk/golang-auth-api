package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func Get() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	})

	http.Handle("/", r)
	return r
}
