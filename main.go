package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	controller "github.com/DantasB/Siga-Professor/Controller"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.HandleFunc("/professors", controller.GetProfessors).Methods("GET")
	router.HandleFunc("/professors/{id}", controller.GetProfessorDisciplines).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))

}
