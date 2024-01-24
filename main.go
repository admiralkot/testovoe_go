package main

import (
	"github.com/gorilla/mux"
	//	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Create a new ServeMux using Gorilla
var rMux = mux.NewRouter()
var log = logrus.New()

func main() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.TextFormatter{})
	// Only log the debug severity or above.
	log.SetLevel(logrus.DebugLevel)
	s := http.Server{
		Addr:         config.PORT,
		Handler:      rMux,
		ErrorLog:     nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	rMux.NotFoundHandler = http.HandlerFunc(DefaultHandler)

	notAllowed := notAllowedHandler{}
	rMux.MethodNotAllowedHandler = notAllowed

	rMux.HandleFunc("/time", TimeHandler)

	// Define Handler Functions
	// Register GET
	getMux := rMux.Methods(http.MethodGet).Subrouter()

	getMux.HandleFunc("/getall", GetAllHandler)
	getMux.HandleFunc("/getid", GetIDHandler)
	getMux.HandleFunc("/username/{id:[0-9]+}", GetUserDataHandler)

	// Register PUT
	// Update User
	putMux := rMux.Methods(http.MethodPut).Subrouter()
	putMux.HandleFunc("/update", UpdateHandler)

	// Register POST
	// Add User
	postMux := rMux.Methods(http.MethodPost).Subrouter()
	postMux.HandleFunc("/add", AddHandler)

	// Register DELETE
	// Delete User
	deleteMux := rMux.Methods(http.MethodDelete).Subrouter()
	deleteMux.HandleFunc("/username/{id:[0-9]+}", DeleteHandler)

	go func() {
		log.Infoln("Listening to", config.PORT)
		err := s.ListenAndServe()
		if err != nil {
			log.Errorln("Error starting server:", err)
			return
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	sig := <-sigs
	log.Infoln("Quitting after signal:", sig)
	time.Sleep(5 * time.Second)
	s.Shutdown(nil)
}
