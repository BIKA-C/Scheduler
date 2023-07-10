package main

import (
	"net/http"
	"scheduler/server"
)

func main() {
	s := http.Server{
		Addr:    ":8080",
		Handler: server.New(),
	}
	s.ListenAndServe()
}
