package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func main() {
	config := GetConfig()

	db, err := GetDBConnection(config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	InitDB(config, db)

	service := NewService(db)
	v := validator.New()
	handler := NewHandler(*service, v)

	r := mux.NewRouter()
	handler.RegisterHandlers(r)

	log.Printf("Start listening port %s", config.ListenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.ListenPort), r))
}
