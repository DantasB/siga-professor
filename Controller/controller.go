package controller

import (
	"fmt"
	"log"
	"net/http"

	crawler "github.com/DantasB/Siga-Professor/Crawler"
	utils "github.com/DantasB/Siga-Professor/Utils"
)

func GetProfessorDisciplines(w http.ResponseWriter, r *http.Request) {
	// Get all professor disciplines
}

func FillDatabase(w http.ResponseWriter, r *http.Request) {
	// Get all siga oppened disciplines
	disciplines, err := crawler.AccessSiraCourses()
	if err != nil {
		log.Fatalln(err)
		return
	}

	//Treat the null disciplines
	fmt.Println("Treating dataset.")
	disciplines, err = utils.FillNilDataWithLastLineData(disciplines)
	if err != nil {
		log.Fatalln(err)
		return
	}

	//Remove the duplicate lines in the disciplines slice
	fmt.Println("Removing Duplicated Lines.")
	disciplines = utils.RemoveDuplicateLines(disciplines)

	//Access the database collections
	disciplinesCollection := utils.GetCollection("disciplines")
	professorsCollection := utils.GetCollection("professors")

	//Generate objects from disciplines slice
	if utils.SaveOnMongo(disciplines, professorsCollection, disciplinesCollection) {
		fmt.Println("Database filled successfully.")
	} else {
		fmt.Println("Error on database filling.")
	}
}
