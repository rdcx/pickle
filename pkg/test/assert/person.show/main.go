package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis/v9"
	"github.com/gorilla/mux"
)

type Person struct {
	ID string `json:"id"`

	Name string `json:"name"`
}

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "redis:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func Function(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	val, err := rdb.Get(ctx, r.URL.Query().Get("id")).Result()
	if err == redis.Nil {
		w.WriteHeader(http.StatusNotFound)
	} else if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		var person Person
		json.Unmarshal([]byte(val), &person)
		json.NewEncoder(w).Encode(person)
	}
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", Function)

	log.Fatal(http.ListenAndServe(":80", r))
}
