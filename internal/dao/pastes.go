package dao

import (
  "context"
  "log"
  "time"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"

  "errors"

  "xyz.haff/pastebox/internal/db"
)

type Paste struct {
  ID string `bson:"_id,omitempty"`
  Title string
  Content string
  Created time.Time
  Expires time.Time `bson:"expires"`
}

type PasteDAO struct {
  collection *mongo.Collection
  infoLog *log.Logger
}

func NewPasteDAO(mongo *mongo.Client, infoLog *log.Logger) *PasteDAO {
  collection := mongo.Database(db.DatabaseName).Collection("pastes")

  return &PasteDAO { collection, infoLog }
}

func (dao *PasteDAO) Insert(title string, content string, expires int) (string, error) {
  result, err := dao.collection.InsertOne(context.TODO(), Paste {
    Title: title,
    Content: content,
    Created: time.Now(),
    Expires: time.Now().AddDate(0, 0, expires),
  } )

  if err != nil {
    return "", err
  }

  id := db.GetInsertOneStringId(result)
  dao.infoLog.Printf("Inserted paste %s", id)

  return id, nil
}

func (dao *PasteDAO) Get(id string) (*Paste, error) {
  objectId, _ := primitive.ObjectIDFromHex(id)

  var result Paste
  err := dao.collection.FindOne(context.TODO(), bson.D{{ "_id", objectId}}).Decode(&result)

  if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments) {
      return nil, db.NewPasteNotFoundError(id)
    } else {
      return nil, err
    }
  }

  return &result, nil
}

// TODO: Handler for this
func (dao *PasteDAO) Latest() ([]Paste, error) {
  opt := options.
          Find().
          SetLimit(10).
          SetSort(bson.D {{ "_id", -1 }})

  nonExpiredFilter := bson.M {
    "expires": bson.M { "$gt": time.Now() },
  }

  cursor, err := dao.collection.Find(context.TODO(), nonExpiredFilter, opt)

  if err != nil {
    return nil, err
  }

  var results []Paste
  if err = cursor.All(context.TODO(), &results); err != nil {
    return nil, err
  }

  return results, nil
}
