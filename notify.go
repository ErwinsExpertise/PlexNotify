package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	h "github.com/ErwinsExpertise/PlexNotify/handlers"

	"github.com/gorilla/mux"
)

var port string

func init() {
	flag.StringVar(&port, "port", "9000", "Port to be used")
	flag.Parse()

	if _, err := os.Stat("/var/log/plexnotify"); os.IsNotExist(err) {
		err := os.Mkdir("/var/log/plexnotify", 0755)
		if err != nil {
			log.Panic(err)
		}
	}

	logFile, err := os.OpenFile("/var/log/plexnotify/notify.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func main() {

	rout := mux.NewRouter()
	rout.HandleFunc("/event", h.EventHandler).Methods("POST")
	rout.HandleFunc("/activity", h.ActivityHandler).Methods("POST", "GET")

	log.Println("Now listening on: " + port)
	log.Println(http.ListenAndServe(":"+port, rout))
}
