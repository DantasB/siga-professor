package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	crawler "github.com/DantasB/Siga-Professor/Crawler"
	models "github.com/DantasB/Siga-Professor/Models"
	utils "github.com/DantasB/Siga-Professor/Utils"
	"github.com/gorilla/mux"
)

func GetProfessorDisciplines(w http.ResponseWriter, r *http.Request) {
	// Get all professor disciplines
	params := mux.Vars(r)
	professor := strings.ToUpper(params["professor"])
	professorObject := models.Professor{}
	professorObject.Name = professor

	disciplinesCollection := utils.GetCollection("disciplines")
	professorsCollection := utils.GetCollection("professors")

	//Get the professor object
	fmt.Println("Trying to find this professor in the database")
	queryResult, found := utils.FoundOne(professorsCollection, utils.BuildTextFilter(professorObject.Name))
	if !found {
		fmt.Println("Professor not found.")
		return
	}

	err := queryResult.Decode(&professorObject)
	if err != nil {
		log.Fatalln(err)
	}

	//Get the disciplines of the professor
	fmt.Println("Trying to find the disciplines of this professor in the database")
	disciplines := utils.GetDisciplinesByProfessor(professorObject, disciplinesCollection)

	//Return the disciplines
	utils.ReturnJSON(w, disciplines)

}

func FillDatabase(w http.ResponseWriter, r *http.Request) {
	// Get all siga oppened disciplines
	disciplines, err := crawler.AccessSiraCourses()
	if err != nil {
		log.Fatalln(err)
	}

	//Treat the null disciplines
	fmt.Println("Treating dataset.")
	disciplines, err = utils.FillNilDataWithLastLineData(disciplines)
	if err != nil {
		log.Fatalln(err)
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
