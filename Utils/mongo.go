package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	models "github.com/DantasB/Siga-Professor/Models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect() *mongo.Client {
	mongoUrl := buildMongoUrl()
	clientOptions := options.Client().ApplyURI(mongoUrl)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalln(err)
	}

	return client
}

func buildMongoUrl() string {
	mongoUrl := os.Getenv("CONNECTION_URL")
	mongoUser := os.Getenv("MONGODB_USER")
	mongoPassword := os.Getenv("MONGODB_PASSWORD")
	mongoDatabase := os.Getenv("MONGODB_DATABASE")
	mongoHost := os.Getenv("MONGODB_HOST")
	mongoParams := os.Getenv("MONGODB_PARAMS")

	mongoUrl = strings.Replace(mongoUrl, "{USER}", mongoUser, 1)
	mongoUrl = strings.Replace(mongoUrl, "{PASSWORD}", mongoPassword, 1)
	mongoUrl = strings.Replace(mongoUrl, "{DATABASE}", mongoDatabase, 1)
	mongoUrl = strings.Replace(mongoUrl, "{HOST}", mongoHost, 1)
	mongoUrl = strings.Replace(mongoUrl, "{PARAMS}", mongoParams, 1)

	return mongoUrl
}

func GetCollection(collection string) *mongo.Collection {
	client := connect()
	return client.Database(os.Getenv("MONGODB_DATABASE")).Collection(collection)
}

func FoundOne(collection *mongo.Collection, filter interface{}) (*mongo.SingleResult, bool) {
	result := collection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return result, false
	}

	return result, true
}

func FindAll(collection *mongo.Collection, filter primitive.M, dataType string) []interface{} {
	var results []interface{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Data not found.")
		return results
	}

	if dataType == "Discipline" {
		FindAllDisciplines(cursor, &results)
	} else if dataType == "Professor" {
		FindAllProfessors(cursor, &results)
	} else {
		log.Fatalln("Object type not found.")
	}

	return results
}

func FindAllProfessors(cursor *mongo.Cursor, professors *[]interface{}) {
	uniqueProfessors := map[string]string{}
	for cursor.Next(context.TODO()) {
		var professor models.Professor
		err := cursor.Decode(&professor)
		if err != nil {
			log.Fatalln(err)
		}
		if _, ok := uniqueProfessors[professor.Name]; ok {
		} else {
			uniqueProfessors[professor.Name] = professor.Name
			*professors = append(*professors, professor)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor
	cursor.Close(context.TODO())
}

func FindAllDisciplines(cursor *mongo.Cursor, disciplines *[]interface{}) {
	uniqueDisciplines := map[string]string{}
	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		discipline := models.Discipline{}
		err := cursor.Decode(&discipline)
		if err != nil {
			log.Fatal(err)
		}
		if _, ok := uniqueDisciplines[discipline.Name]; ok {
		} else {
			uniqueDisciplines[discipline.Name] = discipline.Name
			*disciplines = append(*disciplines, discipline)
		}

	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor
	cursor.Close(context.TODO())
}

func SafeInsertOne(collection *mongo.Collection, document interface{}, filter primitive.M) (*mongo.InsertOneResult, error) {
	_, foundedObject := FoundOne(collection, filter)
	if foundedObject {
		fmt.Println("Object already in the database.")
		return nil, nil
	}

	result, err := InsertOne(collection, document)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Object inserted in the database.")
	return result, err
}

func InsertOne(collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	result, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		log.Fatalln(err)
	}

	return result, err
}

func BuildFilter(document interface{}, objectType string) primitive.M {
	if objectType == "Discipline" {
		return bson.M{
			"nome":         document.(models.Discipline).Name,
			"codigo":       document.(models.Discipline).Code,
			"turma":        document.(models.Discipline).Class,
			"dias":         document.(models.Discipline).Datetime,
			"professor_id": document.(models.Discipline).ProfessorID,
		}
	} else if objectType == "Professor" {
		return bson.M{
			"nome": document.(models.Professor).Name,
		}
	} else {
		log.Fatalln("Object type not found.")
		return nil
	}
}

func BuildTextFilter(text string) primitive.M {
	return bson.M{"$text": bson.M{"$search": text}}
}
