package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"live_server/config"
	"live_server/db"
	"net/http"
	"strconv"
	"time"

	"live_server/settings"
)

func CreateLive(c *gin.Context) {
	name := c.PostForm("name")
	poster := c.PostForm("poster")

	id := settings.GenNewID()
	filter := bson.D{{"name", name}}
	var result settings.Live
	err := db.LiveColl.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			streamID := settings.ServiceName + "/" + strconv.Itoa(id)
			newLive := settings.Live{
				Name:      name,
				Poster:    poster,
				StartTime: time.Now().Round(time.Second),
				RtmpAddr:  config.Config.RtmpPushPullURL + streamID,
				StreamID:  streamID,
			}

			res, er := db.LiveColl.InsertOne(context.TODO(), newLive)
			if er != nil {
				return
			}
			fmt.Printf("Inserted ID: %s\n", res.InsertedID)

			c.JSON(http.StatusOK, newLive)
			fmt.Println(name, poster, streamID)

		} else {
			panic(err)
		}

	}

}

func GetAllLive(c *gin.Context) {
	filter := bson.D{{}}
	sort := bson.D{{"StartTime", 1}}
	cursor, err := db.LiveColl.Find(context.TODO(), filter, options.Find().SetSort(sort))

	// Unpacks the cursor into a slice
	var results []settings.Live
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, results)
}

func fuzzySearchLive(c *gin.Context) {
	name := c.Query("name")
	filter := bson.D{bson.E{Key: "name",
		Value: bson.M{"$regex": primitive.Regex{Pattern: ".*" + name + ".*", Options: "i"}}}}

	cursor, err := db.LiveColl.Find(context.TODO(), filter)

	var results []settings.Live
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, results)
}
