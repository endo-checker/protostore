package store

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
)

type Store[T proto.Message] struct {
	locaColl *mongo.Collection
}

func Connect[T proto.Message](uri, coll string) Store[T] {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("info")
	db.Collection(coll)

	return Store[T]{
		locaColl: db.Collection(coll),
	}
}
