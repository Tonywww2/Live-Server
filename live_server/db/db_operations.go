package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"live_server/settings"
)

func CheckDBContains(coll *mongo.Collection, filter bson.D) bool {
	var result settings.Live
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		return errors.Is(err, mongo.ErrNoDocuments)
	}
	return true

}

func InsertLive(coll *mongo.Collection, live *settings.Live) bool {
	res, er := coll.InsertOne(context.TODO(), live)
	if er != nil {
		return false
	}
	fmt.Printf("Inserted ID: %s\n", res.InsertedID)
	return true
}

func FindLive(coll *mongo.Collection, filter bson.D) (*[]settings.Live, error) {
	cursor, err := coll.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, err
	}
	var results []settings.Live
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}
	return &results, nil

}

func UpdateLive(coll *mongo.Collection, filter bson.D, update bson.D) error {
	res, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Println("Update: ", res)
	return nil
}
