package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting webhooks server...")
	router := mux.NewRouter()

	server := &http.Server{
		Addr:         ":9091",
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	server.ListenAndServe()
}
