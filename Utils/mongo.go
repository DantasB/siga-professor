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

func SafeInsertOne(collection *mongo.Collection, document interface{}, filter primitive.M) (*mongo.InsertOneResult, error) {
	_, foundedObject := FoundOne(collection, filter)
	if foundedObject {
		fmt.Print("Object already in the database.")
		return nil, nil
	}

	result, err := InsertOne(collection, document)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print("Object inserted in the database.")
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
			"name":      document.(models.Discipline).Name,
			"code":      document.(models.Discipline).Code,
			"class":     document.(models.Discipline).Class,
			"datetime":  document.(models.Discipline).Datetime,
			"professor": document.(models.Discipline).ProfessorID,
		}
	} else if objectType == "Professor" {
		return bson.M{
			"name": document.(models.Professor).Name,
		}
	} else {
		log.Fatalln("Object type not found.")
		return nil
	}
}
