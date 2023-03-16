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
	audit      bool

	locaColl *mongo.Collection
}
type StoreOption func(*storeOptions)

type storeOptions struct {
	audit bool
}

func NewStore[T proto.Message](opts ...StoreOption) *Store[T] {
	// apply any options provided
	co := storeOptions{
		audit: true,
	}
	for _, opt := range opts {
		opt(&co)
	}

	// get collection name from proto
	msg := *new(T)
	pbName := string(msg.ProtoReflect().Descriptor().Name())
	pbName = strings.ToLower(pbName)

	return &Store[T]{
		audit: co.audit,

		protoField: pbName,
	}
}

func Connect(uri, coll string) {

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("info")
	db.Collection(coll)
}
