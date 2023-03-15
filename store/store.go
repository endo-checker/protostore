package store

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storer interface {
	Create(ctx context.Context, ent struct{}) error
	Get(ctx context.Context, id string) (struct{}, error)
	Update(id string, ctx context.Context, ent struct{}) error
	Delete(id string) error
}

func (s Store) Create(ctx context.Context, ent struct{}) error {
	_, err := s.locaColl.InsertOne(ctx, ent)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (s Store) Get(ctx context.Context, id string) (struct{}, error) {
	var ent struct{}

	if err := s.locaColl.FindOne(context.Background(), bson.M{"id": id}).Decode(ent); err != nil {
		if err == mongo.ErrNoDocuments {
			return ent, err
		}
		return ent, err
	}

	return ent, nil
}

func (s Store) Update(id string, ctx context.Context, ent struct{}) error {
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
