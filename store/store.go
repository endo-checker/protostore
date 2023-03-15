package store

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	pb "google.golang.org/protobuf/types/known/anypb"
)

type Storer interface {
	Create(ctx context.Context, ent pb.Any) error
	Get(ctx context.Context, id string) (*pb.Any, error)
	Update(id string, ctx context.Context, ent *pb.Any) error
	Delete(id string) error
}

func (s Store) Create(ctx context.Context, ent pb.Any) error {
	_, err := s.locaColl.InsertOne(ctx, ent)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (s Store) Get(ctx context.Context, id string) (*pb.Any, error) {
	var ent *pb.Any

	if err := s.locaColl.FindOne(context.Background(), bson.M{"id": id}).Decode(ent); err != nil {
		if err == mongo.ErrNoDocuments {
			return ent, err
		}
		return ent, err
	}

	return ent, nil
}

func (s Store) Update(id string, ctx context.Context, ent *pb.Any) error {
	insertResult, err := s.locaColl.ReplaceOne(ctx, bson.M{"id": id}, ent)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nInserted a Single Document: %v\n", insertResult)

	return err
}

func (s Store) Delete(id string) error {
	if _, err := s.locaColl.DeleteOne(context.Background(), bson.M{"id": id}); err != nil {
		return err
	}
	return nil
}
