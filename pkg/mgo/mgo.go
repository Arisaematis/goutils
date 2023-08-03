package mgo

import (
	"goutils/pkg/drivers"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewMgo(collection string) *Mgo {
	return &Mgo{Collection: drivers.DefaultMongoDatabase.Collection(collection)}
}

type Mgo struct {
	*mongo.Collection
	Pip mongo.Pipeline
}
