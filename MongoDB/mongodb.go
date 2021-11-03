package mongodb

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

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

func GetCollection() *mongo.Collection {
	client := connect()
	mongoDatabase := os.Getenv("MONGODB_DATABASE")
	mongoCollection := os.Getenv("MONGODB_COLLECTION")
	collection := client.Database(mongoDatabase).Collection(mongoCollection)

	return collection
}
