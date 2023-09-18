package models

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"xyz.haff/pastebox/internal/db"
)

type User struct {
  ID string `bson:"_id,omitempty"`
  Name string `bson:"name"`
  Email string `bson:"email"`
  HashedPassword []byte `bson:"hashedPassword"`
  Created time.Time `bson:"created"`
}

type UserDAO struct {
  collection *mongo.Collection
  infoLog *log.Logger
}

func NewUserDAO(mongo *mongo.Client, infoLog *log.Logger) *UserDAO {
  collection := mongo.Database(db.DatabaseName).Collection("users")

  return &UserDAO { collection, infoLog }
}

func (dao *UserDAO) Insert(name, email, password string) error {
  //user := User { 
    //Name: name,
    //Email: email,
    //Created: time.Now().Truncate(time.Second),
  //}

  return nil
}

func (dao *UserDAO) Authenticate(email, password string) (int, error) {
  return 0, nil
}

func (dao *UserDAO) Exists(id string) (bool, error) {
  return false, nil
}
