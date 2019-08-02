package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// SushiAPI holds shared dependencies
type SushiAPI struct {
	dbClient *mongo.Client
	router   *mux.Router
}

// Start the SushiAPI.
func (s *SushiAPI) Start(host string) {
	s.init()
	log.Fatal(http.ListenAndServe(host, s.router))
	s.shutdown()
}

func (s *SushiAPI) init() {
	s.router.HandleFunc("/sushi", s.handleGetRolls()).Methods("GET")
	s.router.HandleFunc("/sushi/{id}", s.handleGetRoll()).Methods("GET")
	s.router.HandleFunc("/sushi", s.handleCreateRoll()).Methods("POST")
	s.router.HandleFunc("/sushi/{id}", s.handleUpdateRoll()).Methods("PUT")
	s.router.HandleFunc("/sushi/{id}", s.handleDeleteRoll()).Methods("DELETE")
}

// Shutdown the SushiAPI.
func (s *SushiAPI) shutdown() {
	err := s.dbClient.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connect to MongoDB closed.")
}
