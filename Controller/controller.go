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
	fmt.Println("Building Objects.")
	disciplineObjects, professorObjects := utils.ConvertToObjects(disciplines, professorsCollection, disciplinesCollection)

	for _, professor := range professorObjects {
		filter := utils.BuildFilter(professor, "professor")
		fmt.Printf("Trying to insert the Professor %v in the database.\n", professor.Name)
		utils.SafeInsertOne(professorsCollection, professor, filter)
	}

	for _, discipline := range disciplineObjects {
		filter := utils.BuildFilter(discipline, "discipline")
		fmt.Printf("Trying to insert the Discipline %v in the database.\n", discipline.Name)
		utils.SafeInsertOne(disciplinesCollection, discipline, filter)
	}

}
