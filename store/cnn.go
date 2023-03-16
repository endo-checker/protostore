package store

import (
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
)

type Store[T proto.Message] struct {
	protoField string
	locaColl   *mongo.Collection
}

// add your mongo uri, and collection name
// connect to your proto.Message type
// e.g. store.Connect[*proto.Message]("mongodb://localhost:27017", "info")
func Connect[T proto.Message](uri, coll string) Store[T] {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("info")
	db.Collection(coll)

	msg := *new(T)
	pbName := string(msg.ProtoReflect().Descriptor().Name())
	pbName = strings.ToLower(pbName)

	return Store[T]{
		locaColl:   db.Collection(coll),
		protoField: pbName,
	}
}
