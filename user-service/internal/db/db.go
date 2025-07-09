package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var  (

	DbConnectTimeout time.Duration = 10*time.Second

	Client *mongo.Client = NewDbClient()

)



func NewDbClient() *mongo.Client {
	uri := os.Getenv("URI")

	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: true,
		NilMapAsEmpty:     true,
		NilSliceAsEmpty:   true,
	}
	
	clientOpts := options.Client().ApplyURI(uri).SetBSONOptions(bsonOpts)

	ctx,cancel := context.WithTimeout(context.Background(),DbConnectTimeout);
	defer cancel()

	client , err := mongo.Connect(ctx,clientOpts)
	if err!=nil {
		log.Panicln("database connection failed: ", err)
		return nil
	}

	return client
}

