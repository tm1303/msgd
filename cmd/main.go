package main

import (
	"context"
	_ "embed"
	"fmt"
	"msgd/broadcaster"
	"msgd/domain"
	"msgd/infra"
	"msgd/processor"
	"msgd/receiver"
	"msgd/ui"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-chi/chi/middleware"
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
	pollSize:  10,
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All good\n")) //nolint:errcheck
}

func main() {
	ctx := context.Background()

	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(defaultConfig.awsRegion),
	}))

	broadcastChan := make(chan domain.MessageBody, 100)
	msgPoller := processor.NewSqsPoller(awsSession, defaultConfig.queueName, defaultConfig.pollSize)
	go broadcaster.StartBroadcaster(ctx, broadcastChan)
	go processor.StartProcessor(ctx, msgPoller, broadcastChan)

	msgClient := receiver.NewSqsQueuer(awsSession, defaultConfig.queueName)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
	r.Use(infra.UserIDMiddleware)

	r.Get("/", ui.ServeHTML)

	r.Get("/health", health)
	r.Post("/enqueue", receiver.GetHandler(msgClient))
	r.Get("/ws", broadcaster.WsHandler)

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

	ctx.Done() // stops anything listening to our main context, ie the processor

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // new context with timeout to control shutdown
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server %s\n", err)
	}
	fmt.Println("Server closed cleanly")
}
