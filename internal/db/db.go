package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: Pass it from somewhere?
const mongoUri = "mongodb://root:pass@localhost:27017"
func ConnectToMongo(infoLogger *log.Logger) (*mongo.Client, func(), error) {
  serverAPI := options.ServerAPI(options.ServerAPIVersion1)

  opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)

  timeoutCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()
  
  client, err := mongo.Connect(timeoutCtx, opts)
  if err != nil {
    return nil, nil, err
  }
  close := func() {
    timeoutCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()
    if err = client.Disconnect(timeoutCtx); err != nil {
      panic(err)
    }
  }

  var result bson.M
  if err := client.Database("admin").RunCommand(context.TODO(), bson.M {"ping": 1}).Decode(&result); err != nil {
	  close()	
    return nil, nil, err
	}
  // TODO: ooops, logging the password
  infoLogger.Printf("Connected to mongo on %s", mongoUri)
  return client, close, nil
}
