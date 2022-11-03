package main

import (
	"context"
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
	err := rdb.Del(ctx, "person:"+r.URL.Query().Get("id")).Err()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", Function)

	log.Fatal(http.ListenAndServe(":80", r))
}
