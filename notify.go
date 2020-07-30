package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ErwinsExpertise/PlexNotify/handlers"
)

var port string

const (
	POST = "POST"
	GET  = "GET"
)

func init() {
	flag.StringVar(&port, "port", "9000", "Port to be used")
	flag.Parse()

	var logPath string

	if _, err := os.Stat("/var/log/plexnotify"); os.IsNotExist(err) {
		err := os.Mkdir("/var/log/plexnotify", 0755)
		if err != nil {
			logPath = ""
		} else {
			logPath = "/var/log/plexnotify/"
		}
	}

	logFile, err := os.OpenFile(logPath+"notify.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	// Create channel of size 1
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			sig := <-c
			log.Printf("Recieved %s signal. Shutting down server...\n", sig)
			cancel()
			os.Exit(0)
		case <-ctx.Done():
			// consume
		}

	}()

	rout := mux.NewRouter()

	rout.HandleFunc("/event", handlers.EventHandler).Methods(POST)
	rout.HandleFunc("/activity", handlers.ActivityHandler).Methods(POST, GET)
	rout.HandleFunc("/logs", handlers.LogHandler).Methods(POST, GET)
	rout.HandleFunc("/login", handlers.LoginHandler).Methods(POST, GET)
	rout.HandleFunc("/", handlers.ActivityHandler).Methods(POST, GET)

	// route for Prometheus metrics
	rout.Handle("/metrics", promhttp.Handler())

	log.Println("Now listening on: " + port)
	log.Println(http.ListenAndServe(":"+port, rout))
}
