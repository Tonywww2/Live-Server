package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"live_server/settings"
	"net/http"
	"strconv"
	"time"
)

func CreateLive(c *gin.Context) {
	name := c.PostForm("name")
	poster := c.PostForm("poster")

	id := settings.GenNewID()

	flag := true
	for _, v := range settings.LiveList {
		if v.Name == name {
			c.String(http.StatusNotAcceptable, fmt.Sprintf("Already exist"))
			flag = false
		}
	}
	if flag {
		streamID := settings.ServiceName + "/" + strconv.Itoa(id)
		settings.LiveList[streamID] = settings.Live{
			name,
			poster,
			time.Now().Round(time.Second),
			settings.RtmpPushPullURL + streamID,
			streamID,
			false,
		}

		c.JSON(http.StatusOK, settings.LiveList[streamID])
		fmt.Println(name, poster, streamID)

		result, err := settings.Coll.InsertOne(context.TODO(), settings.LiveList[streamID])
		if err != nil {
			return
		}
		fmt.Println("Inserted ID: %s", result.InsertedID)

	}

}

func GetLiveName(c *gin.Context) {
	id := c.DefaultQuery("stream_id", "")
	fmt.Println(id)
	if result, contains := settings.LiveList[id]; contains {
		c.String(http.StatusOK, result.ToString())

	} else {
		c.String(http.StatusNotFound, "")

	}

}

func GetAllLive(c *gin.Context) {
	c.JSON(http.StatusOK, settings.LiveList)
}
