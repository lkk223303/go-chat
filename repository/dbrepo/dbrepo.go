package dbrepo

import (
	"chatty/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDBRepo struct {
	Client *mongo.Client
}

func NewMongoRepo(c *mongo.Client) repository.DatabaseRepo {
	return &mongoDBRepo{
		Client: c,
	}
}
