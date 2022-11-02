package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Function(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode("Hello world")
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", Function).Methods("GET")

	log.Fatal(http.ListenAndServe(":80", r))
}
