package store

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	locaColl *mongo.Collection
}

func Connect(mongoUri string, dbGroup string, dbColl string) Store {
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(dbGroup)

	return Store{
		locaColl: db.Collection(dbColl),
	}
}
