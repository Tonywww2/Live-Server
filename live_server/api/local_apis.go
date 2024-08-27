package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"live_server/config"
	"live_server/db"
	"net/http"
	"os"
	"strconv"
	"time"

	"live_server/settings"
)

type LiveApi struct {
	LiveColl *mongo.Collection
}

// 这里是构造函数
func NewLiveApi() *LiveApi {
	return &LiveApi{
		LiveColl: db.GetCollection("live_list"),
	}
}

func (this *LiveApi) RegisterRouter(server *GinServer) {
	// RegisterRouter /******************Start 注册直播间路由*******************/
	server.POST("/createLive", this.CreateLive)
	server.GET("/getAllLive", this.GetAllLive)
	server.GET("/fuzzySearchLive", this.FuzzySearchLive)
	server.GET("/getRecordList", this.GetRecordList)
	// RegisterRouter /******************End 注册直播间路由*******************/
}

func (this *LiveApi) CreateLive(c *gin.Context) {
	name := c.PostForm("name")
	poster := c.PostForm("poster")
	id := settings.GenNewID()
	filter := bson.D{{"name", name}}

	if !db.CheckDBContains(this.LiveColl, filter) {
		streamID := settings.ServiceName + "/" + strconv.Itoa(id)
		newLive := settings.Live{
			Name:      name,
			Poster:    poster,
			StartTime: time.Now().Round(time.Second),
			RtmpAddr:  config.Config.RtmpPushPullURL + streamID,
			StreamID:  streamID,
		}

		if db.InsertLive(this.LiveColl, &newLive) {
			c.JSON(http.StatusOK, newLive)
			fmt.Println(name, poster, streamID)

		} else {
			c.JSON(http.StatusNotAcceptable, "DB Error")
		}

	} else {
		c.JSON(http.StatusConflict, "Live already found")
	}

}

func (this *LiveApi) GetAllLive(c *gin.Context) {
	filter := bson.D{{}}
	sort := bson.D{{"StartTime", 1}}
	cursor, err := this.LiveColl.Find(context.TODO(), filter, options.Find().SetSort(sort))

	// Unpacks the cursor into a slice
	var results []settings.Live
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, results)
}

func (this *LiveApi) FuzzySearchLive(c *gin.Context) {
	name := c.Query("name")
	filter := bson.D{bson.E{Key: "name",
		Value: bson.M{"$regex": primitive.Regex{Pattern: ".*" + name + ".*", Options: "i"}}}}

	res, err := db.FindLive(this.LiveColl, filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (this *LiveApi) GetRecordList(c *gin.Context) {
	var files []string
	files, err := this.GetAllFile(config.Config.M7sRecordDir, files)
	fmt.Println(files)
	if err != nil {
		fmt.Println("Error reading the dir: ", err)
		c.JSON(http.StatusInternalServerError, err)

	} else {
		c.JSON(http.StatusOK, files)
	}
}

func (this *LiveApi) GetAllFile(pathname string, s []string) ([]string, error) {
	rd, err := os.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = this.GetAllFile(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}
