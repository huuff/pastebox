package models

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"xyz.haff/pastebox/internal/db"
  customError "xyz.haff/pastebox/internal/errors"
)

type User struct {
  ID string `bson:"_id,omitempty"`
  Name string `bson:"name"`
  Email string `bson:"email"`
  HashedPassword []byte `bson:"hashedPassword"`
  Created time.Time `bson:"created"`
}

type UserDAOInterface interface {
  Insert(name, email, password string, ctx context.Context) (string, error)
  Authenticate(email, password string, ctx context.Context) (string, error)
  Exists(id string, ctx context.Context) (bool, error)
}

type UserDAO struct {
  collection *mongo.Collection
  infoLog *log.Logger
}

func NewUserDAO(client *mongo.Client, infoLog *log.Logger, ctx context.Context) *UserDAO {
  collection := client.Database(db.DatabaseName).Collection("users")

  _, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel {
    Keys: bson.M {"email": 1},
    Options: options.Index().SetUnique(true),
  })

  if err != nil {
    panic(err)
  }

  return &UserDAO { collection, infoLog }
}

func (dao *UserDAO) Insert(name, email, password string, ctx context.Context) (string, error) {
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

  if err != nil {
    return "", err
  }

  user := User { 
    Name: name,
    Email: email,
    Created: time.Now().Truncate(time.Second),
    HashedPassword: hashedPassword,
  }

  result, err := dao.collection.InsertOne(ctx, user )

  if err != nil {
    if mongo.IsDuplicateKeyError(err) {
      return "", customError.DuplicateEmailError
    }
    return "", err
  }

  id := db.GetInsertOneStringId(result)
  dao.infoLog.Printf("Inserted user %s", id)

  return "", nil
}

func (dao *UserDAO) Authenticate(email, password string, ctx context.Context) (string, error) {
  var user User
  err := dao.collection.FindOne(ctx, bson.M {
    "email": email,
  }).Decode(&user)

  if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments) {
      return "", db.ErrInvalidCredentials
    } else {
      return "", err
    }
  }

  if err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)); err != nil {
    if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
      return "", db.ErrInvalidCredentials
    } else {
      return "", err
    }
  }

  return user.ID, nil
}

func (dao *UserDAO) Exists(id string, ctx context.Context) (bool, error) {
  objectId, err := primitive.ObjectIDFromHex(id)
  if err != nil {
    return false, err
  }

  var user User 
  err = dao.collection.FindOne(ctx, bson.M { "_id": objectId}).Decode(&user)

  if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments) {
      return false, nil
    } else {
      return false, err
    }
  }

  return true, nil
}
