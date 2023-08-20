package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
  "fmt"
)

func GetInsertOneStringId(r *mongo.InsertOneResult) string {
  objectId, ok := r.InsertedID.(primitive.ObjectID)

  if !ok {
    panic(fmt.Sprintf("Cannot convert %v to ObjectId", r.InsertedID))  
  }

  return objectId.String()
}

