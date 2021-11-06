package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	controller "github.com/DantasB/Siga-Professor/Controller"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/filldatabase", controller.FillDatabase).Methods("GET")
	router.HandleFunc("/professors/{id}", controller.GetProfessorDisciplines).Methods("GET")

	http.ListenAndServe(":8000", router)

}
