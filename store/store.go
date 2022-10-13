package store

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	coll *mongo.Collection
}

// retrives environment variables
func LoadEnv(env string) (uri string) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file ")
	}
	uri = os.Getenv(env)
	return uri
}

func Connect(i, dbName, collName string) Store {
	uri := LoadEnv(i)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(dbName)

	return Store{
		coll: db.Collection(collName),
	}
}
