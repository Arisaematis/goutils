package drivers

import (
	"context"
	"goutils/pkg/setting"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	DefaultMongoClient   *mongo.Client
	DefaultMongoDatabase *mongo.Database
)

func init() {
	connectCtx, connectCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer connectCancel()
	client, err := mongo.Connect(connectCtx,
		options.Client().
			ApplyURI(setting.RemoteSetting.MongoAddr))
	if err != nil {
		panic(err)
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	DefaultMongoDatabase = client.Database(setting.RemoteSetting.MongoDatabase)
	DefaultMongoClient = client
}
