package controller

import (
	"net/http"

	crawler "github.com/DantasB/Siga-Professor/Crawler"
)

func GetProfessor(w http.ResponseWriter, r *http.Request) {
	// Get all professor disciplines
	return
}

func GetProfessors(w http.ResponseWriter, r *http.Request) {
	// Get all professors
	crawler.AccessSiraCourses()
	return
}
