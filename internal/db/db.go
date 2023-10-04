package db

import (
	"context"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: Pass it from somewhere?
const mongoUri = "mongodb://root:pass@localhost:27017"
func ConnectToMongo(infoLogger *log.Logger, ctx context.Context) (*mongo.Client, func(), error) {
  serverAPI := options.ServerAPI(options.ServerAPIVersion1)

  opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)

  client, err := mongo.Connect(ctx, opts)
  if err != nil {
    return nil, nil, err
  }
  close := func() {
    if err = client.Disconnect(ctx); err != nil {
      panic(err)
    }
  }

  var result bson.M
  if err := client.Database("admin").RunCommand(ctx, bson.M {"ping": 1}).Decode(&result); err != nil {
	  close()	
    return nil, nil, err
	}
  // TODO: ooops, logging the password
  infoLogger.Printf("Connected to mongo on %s", mongoUri)
  return client, close, nil
}
