package main

import (
	"context"
	"fmt"
	"msgd/processor"
	"msgd/receiver"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-chi/chi/v5"
)

type config struct {
	httpPort  string
	queueName string
	awsRegion string
	pollSize  int64
}

var defaultConfig = config{
	httpPort:  ":8080",
	queueName: "https://sqs.eu-west-2.amazonaws.com/205962165374/MSGC-DEV-1",
	awsRegion: "eu-west-2",
	pollSize:  1,
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All good\n"))
}

func main() {
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(defaultConfig.awsRegion), // Replace with your desired region
	}))

	msgPoller := processor.NewSqsPoller(awsSession, defaultConfig.queueName, defaultConfig.pollSize)
	processor.StartProcessor(msgPoller)
	
	msgClient := receiver.NewSqsClient(awsSession, defaultConfig.queueName)

	r := chi.NewRouter()
	r.Get("/health", health)
	r.Post("/enqueue", receiver.GetHandler(msgClient))

	srv := http.Server{
		Addr:    defaultConfig.httpPort,
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("Starting server on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Listen and Serve err %s\n", err)
		}
	}()

	sig := <-stop
	fmt.Printf("\nsignal: %s\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server %s\n", err)
	}
	fmt.Println("Server closed cleanly")
}
