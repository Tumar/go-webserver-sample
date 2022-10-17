package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("POSTGRES_HOST")
	port     = os.Getenv("POSTGRES_PORT")
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

func initDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeoutExceeded := time.After(15 * time.Second)

	for {
		select {
		case <-timeoutExceeded:
			return nil, errors.New("couldn't connect to databse")
		case <-ticker.C:
			err = db.Ping()
			if err == nil {
				return db, err
			}
			fmt.Println("ping failed")
		}
	}
}
