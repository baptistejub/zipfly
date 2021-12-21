package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	zipfly "github.com/baptistejub/zipfly/zip_fly"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "6969"
	}

	publicUrl := os.Getenv("PUBLIC_URL")
	if publicUrl == "" {
		host, err := os.Hostname()
		if err != nil {
			panic(err)
		}

		publicUrl = "http://" + host + ":" + port
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	validateSignature := os.Getenv("VALIDATE_SIGNATURE") != ""
	options := zipfly.ServerOptions{
		ValidateSignature: validateSignature,
		SigningSecret:     os.Getenv("SIGNING_SECRET"),
		PublicUrl:         publicUrl,
	}

	httpServer := &http.Server{
		Addr:        ":" + port,
		Handler:     zipfly.NewServer(environment, options),
		ReadTimeout: 10 * time.Second,
	}

	go func() {
		httpServer.ListenAndServe()
	}()

	log.Printf("Server started on port %s", port)

	// Gracefully shutdown when SIGTERM is received
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Printf("Shutting down...")
	httpServer.Shutdown(context.Background())
}
