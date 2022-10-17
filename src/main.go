package main

import (
	"context"
	"database/sql"
	"errors"
	"gowebserver/src/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "image/png"
)

func startServer(db *sql.DB) *http.Server {
	minioClient, err := newMinioClient()
	if err != nil {
		panic(err)
	}

	r := routes.NewRouter(&routes.Env{Db: db, MinioClient: minioClient})
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	return srv
}

func main() {
	db, err := initDB()
	if err != nil {
		panic(err)
	}

	srv := startServer(db)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer db.Close()
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
