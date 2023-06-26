package main

import (
	"net/http"
	"scheduler/api"
)

func main() {
	s := http.Server{
		Addr:    ":8080",
		Handler: api.New(),
	}
	s.ListenAndServe()
}
