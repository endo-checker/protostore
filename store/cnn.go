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
	locaColl   *mongo.Collection
	protoField string
}

// add your mongo uri, and collection name
// connect to your proto.Message type
// e.g. store.Connect[*proto.Message]("mongodb://localhost:27017", "info")
func Connect[T proto.Message](uri string, opts ...ClientOption) error {

	var err error

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	msg := *new(T)
	pbName := string(msg.ProtoReflect().Descriptor().Name())
	pbName = strings.ToLower(pbName)

	db := client.Database("info")
	db.Collection(pbName)

	return nil
}
