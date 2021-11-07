package utils

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	models "github.com/DantasB/Siga-Professor/Models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var layout = "02/01/2006"

func ToUtf8(iso8859_1_buf []byte) string {
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}

func ParseDate(date string) (time.Time, error) {
	splittedDate := strings.Split(date, " ")
	if len(splittedDate) != 4 {
		return time.Time{}, errors.New("could not split the string")
	}

	parsedTime, err := time.Parse(layout, splittedDate[2])
	return parsedTime, err
}

func RemoveSeparators(data string) string {
	return strings.Replace(strings.Replace(strings.Replace(data, "\t", " ", -1), "\r", " ", -1), "\n", " ", -1)
}

func ReplaceMultipleSpacesByPipe(data string) string {
	regex := regexp.MustCompile(` {6,}`)
	return regex.ReplaceAllString(data, "|")
}

func FillNilDataWithLastLineData(disciplines [][]string) ([][]string, error) {
	lastLine := []string{}
	for _, line := range disciplines {
		if len(line) == 0 {
			continue
		}

		if line[0] != "\u00a0" {
			lastLine = line
		}

		for columnIndex, column := range line {
			if column == "\u00a0" {
				line[columnIndex] = lastLine[columnIndex]
			}
		}
	}

	return disciplines, nil
}

func RemoveDuplicateLines(disciplines [][]string) [][]string {
	type data struct {
		code      string
		class     string
		name      string
		day       string
		time      string
		professor string
	}

	deduplicatedDisciplines := [][]string{}
	duplicates := map[data]int{}
	for index, line := range disciplines {
		if len(line) == 0 {
			continue
		}

		information := data{line[0], line[1], line[2], line[3], line[4], line[5]}
		if _, ok := duplicates[information]; ok {
		} else {
			duplicates[information] = index
			deduplicatedDisciplines = append(deduplicatedDisciplines, line)
		}
	}

	return deduplicatedDisciplines
}

func ConvertToObjects(disciplines [][]string, professorsCollection *mongo.Collection, disciplinesCollection *mongo.Collection) ([]models.Discipline, []models.Professor) {
	disciplinesObjects := []models.Discipline{}
	professorsObjects := []models.Professor{}
	for _, line := range disciplines {

		professorObject := generateProfessorObject(professorsCollection, line, &professorsObjects)

		generateDisciplineObject(disciplinesCollection, line, professorObject.ID, &disciplinesObjects)
	}

	return disciplinesObjects, professorsObjects
}

func generateDisciplineObject(collection *mongo.Collection, line []string, professorID primitive.ObjectID, disciplinesObjects *[]models.Discipline) {
	filter := bson.M{
		"name":      line[2],
		"code":      line[0],
		"class":     line[1],
		"datetime":  models.Datetime{line[3], line[4]},
		"professor": professorID,
	}

	queryResult, found := FoundOne(collection, filter)
	if found {
		discipline := models.Professor{}
		err := queryResult.Decode(discipline)
		if err != nil {
			log.Fatalln(err)
		}

		return
	}

	disciplineObject := models.Discipline{
		ID:          primitive.NewObjectID(),
		Name:        line[2],
		Code:        line[0],
		Class:       line[1],
		Datetime:    models.Datetime{line[3], line[4]},
		ProfessorID: professorID,
	}

	*disciplinesObjects = append(*disciplinesObjects, disciplineObject)

}

func generateProfessorObject(collection *mongo.Collection, line []string, professorsObjects *[]models.Professor) models.Professor {

	filter := bson.M{"name": line[5]}

	queryResult, found := FoundOne(collection, filter)
	if found {
		professor := models.Professor{}
		err := queryResult.Decode(professor)
		if err != nil {
			log.Fatalln(err)
		}

		return professor
	}

	professorObject := models.Professor{
		ID:   primitive.NewObjectID(),
		Name: line[5],
	}

	*professorsObjects = append(*professorsObjects, professorObject)

	_, err := InsertOne(collection, professorObject)
	if err != nil {
		log.Fatalln(err)
	}

	return professorObject
}
