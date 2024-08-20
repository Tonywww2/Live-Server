package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"live_server/api"
	"live_server/settings"
)

func main() {
	// DB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(settings.MongodbUri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	settings.Coll = client.Database("live_list").Collection("live_list")

	// Init Map
	filter := bson.D{{}}
	cursor, err := settings.Coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	// Unpacks the cursor into a slice
	var result []settings.Live
	if err = cursor.All(context.TODO(), &result); err != nil {
		panic(err)
	}
	for _, v := range result {
		settings.LiveList[v.StreamID] = v
	}
	fmt.Println(settings.LiveList)

	r := gin.Default()
	r.POST("/createLive", api.CreateLive)
	r.GET("/getLiveName", api.GetLiveName)
	r.GET("/getAllLive", api.GetAllLive)
	r.POST("/pushVideoToStream", api.PushVideoToStream)
	r.POST("/pushStreamToRtmp", api.PushStreamToRtmp)
	r.POST("/endStream", api.EndStreamAPI)

	r.Run(":8082")
}
