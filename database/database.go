package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"os"
)

var DbConn *sql.DB

var host = os.Getenv("HOST_NAME")

const (
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "db"
)

// SetupDatabase ...
func SetupDatabase() {
	var err error
	if len(host) < 1 {
		host = "localhost"
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	DbConn, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Print(err)
	}
	log.Printf("Successfully connected to database using gorm from %s:%d", host, port)

	DbConn.SetMaxOpenConns(3)
	DbConn.SetMaxIdleConns(3)
	DbConn.SetConnMaxLifetime(60 * time.Second)
}
