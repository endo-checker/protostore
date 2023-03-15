package store

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.starlark.net/lib/proto"
	pb "google.golang.org/protobuf/types/known/anypb"
)

type Storer[T proto.Message] interface {
	Create(ctx context.Context, msg T) error
	Query(ctx context.Context, qr ...listOption) ([]T, int64, error)
	Get(ctx context.Context, id string) (*T, error)
	Update(id string, ctx context.Context, u *T) error
	Delete(id string) error
}

func (s Store[T]) Create(ctx context.Context, msg T) error {
	_, err := s.locaColl.InsertOne(ctx, msg)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

type listOptions struct {
	findOpts options.FindOptions
	filter   bson.M
}

// type ListOption func(*listOptions)
type ListOption interface {
	apply(*listOptions)
}

type listOption struct {
	applyFunc func(*listOptions)
}

func (l *listOption) apply(lo *listOptions) {
	l.applyFunc(lo)
}

// List returns a list of documents matching the filter provided.
func (s Store[T]) List(ctx context.Context, opts ...ListOption) ([]T, int64, error) {
	lo := listOptions{}
	for _, opt := range opts {
		opt.apply(&lo)
	}

	if lo.findOpts.Limit == nil || *lo.findOpts.Limit == 0 {
		var lim int64 = 50
		lo.findOpts.Limit = &lim
	}

	cursor, err := s.locaColl.Find(ctx, &lo.findOpts)
	if err != nil {
		return nil, 0, err
	}

	// unpack results
	var docs []T
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, err
	}

	// count of all matching docs
	matches, err := s.locaColl.CountDocuments(ctx, lo.filter)
	if err != nil {
		return nil, 0, err
	}

	return docs, matches, nil
}

func (s Store[T]) Get(ctx context.Context, id string) (*pb.Any, error) {
	var msg *pb.Any

	if err := s.locaColl.FindOne(context.Background(), bson.M{"id": id}).Decode(msg); err != nil {
		if err == mongo.ErrNoDocuments {
			return msg, err
		}
		return msg, err
	}

	return msg, nil
}

func (s Store[T]) Update(id string, ctx context.Context, u *T) error {
	insertResult, err := s.locaColl.ReplaceOne(ctx, bson.M{"id": id}, u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nInserted a Single Document: %v\n", insertResult)

	return err
}

func (s Store[T]) Delete(id string) error {
	if _, err := s.locaColl.DeleteOne(context.Background(), bson.M{"id": id}); err != nil {
		return err
	}
	return nil
}
