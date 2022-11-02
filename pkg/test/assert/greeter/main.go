package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Function(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode("Hello world")

}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", Function)

	log.Fatal(http.ListenAndServe(":80", r))
}
