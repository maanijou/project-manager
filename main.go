package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/maanijou/project-manager/database"
	"github.com/maanijou/project-manager/employee"
	"github.com/maanijou/project-manager/project"
)

func main() {
	log.Println("Running Project Manager app.")
	log.Println("Setup database...")
	database.SetupDatabase()
	defer database.DbConn.Close()

	sm := mux.NewRouter().StrictSlash(true) // ignoring trailing slash
	sm = sm.PathPrefix("/api/v1/").Subrouter()
	sm.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})
	employee.SetupRoutes(sm)
	project.SetupRoutes(sm)

	s := http.Server{
		Addr:         ":8080",           // configure the bind address
		Handler:      sm,                // set the default handler
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		log.Println("Starting server on port 8080")

		err := s.ListenAndServe()
		if err != nil {
			log.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
