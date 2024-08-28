package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"live_server/settings"
)

type LiveCollStruct struct {
	*mongo.Collection
}

func (coll *LiveCollStruct) CheckDBContains(filter bson.D) bool {
	var result settings.Live
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false
		}
		panic(err)
	}
	return true

}

func (coll *LiveCollStruct) InsertLive(live *settings.Live) bool {
	res, er := coll.InsertOne(context.TODO(), live)
	if er != nil {
		return false
	}
	fmt.Printf("Inserted ID: %s\n", res.InsertedID)
	return true
}

func (coll *LiveCollStruct) FindLive(filter bson.D, sort bson.D) (*[]settings.Live, error) {
	cursor, err := coll.Find(context.TODO(), filter, options.Find().SetSort(sort))
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

func (coll *LiveCollStruct) UpdateLive(filter bson.D, update bson.D) error {
	res, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Println("Update: ", res)
	return nil
}
