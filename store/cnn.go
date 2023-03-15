package store

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.starlark.net/lib/proto"
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
