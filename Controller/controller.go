package controller

import (
	"log"
	"net/http"

	crawler "github.com/DantasB/Siga-Professor/Crawler"
	utils "github.com/DantasB/Siga-Professor/Utils"
)

func GetProfessorDisciplines(w http.ResponseWriter, r *http.Request) {
	// Get all professor disciplines
}

func GetProfessors(w http.ResponseWriter, r *http.Request) {
	// Get all professors
	disciplines, err := crawler.AccessSiraCourses()
	if err != nil {
		log.Fatalln(err)
		return
	}

	disciplines, err = utils.FillNilDataWithLastLineData(disciplines)
	if err != nil {
		log.Fatalln(err)
		return
	}

}
