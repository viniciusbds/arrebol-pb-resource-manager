package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/logger"

	"github.com/joho/godotenv"
	"github.com/viniciusbds/arrebol-pb-resource-manager/api"
	"github.com/viniciusbds/arrebol-pb-resource-manager/balancer"
	"github.com/viniciusbds/arrebol-pb-resource-manager/storage"
)

const (
	ServerPort = "SERVER_PORT"
)

func main() {

	var wait time.Duration
	flag.DurationVar(&wait, "graceful_timeout", time.Second*15, "the duration for which the server "+
		"gracefully wait for existing connections to finish - e.g. 15s or 1m")

	err := godotenv.Load()
	if err != nil {
		logger.Errorln(err.Error())
	}

	apiPort := flag.String(ServerPort, os.Getenv("DEFAULT_SERVER_PORT"), "Service port")

	flag.Parse()

	s := storage.New(os.Getenv("DATABASE_ADDRESS"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_PASSWORD"))
	s.Setup()
	defer s.Driver().Close()

	a := api.New(s)

	// Shutdown gracefully
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		<-sigs
		log.Println("Shutting down service")

		if err := a.Shutdown(); err != nil {
			log.Printf("Failed to shutdown the server: %v", err)
		}
	}()

	go balancer.Start()

	if err := a.Start(*apiPort); err != nil {
		log.Fatal(err.Error())
	}
}