package store

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
)

type Store[T proto.Message] struct {
	protoField string
	audit      bool

	locaColl *mongo.Collection
}

func NewStore[T proto.Message](protoField string, audit bool, locaColl *mongo.Collection) Store[T] {
	return Store[T]{
		protoField: protoField,
		audit:      audit,
		locaColl:   locaColl,
	}
}

func Connect(uri, coll string) error {

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("info")
	db.Collection(coll)

	return nil
}
