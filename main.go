package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Running Project Manager app.")

	sm := mux.NewRouter().StrictSlash(true) // ignoring trailing slash
	sm.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	s := http.Server{
		Addr:         ":8080",           // configure the bind address
		Handler:      sm,                // set the default handler
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
