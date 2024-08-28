package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"strconv"
	"time"

	"live_server/config"
	"live_server/db"
	"live_server/settings"
)

type LiveApi struct {
	LiveColl *db.LiveCollStruct
}

// 这里是构造函数
func NewLiveApi() *LiveApi {
	return &LiveApi{
		&db.LiveCollStruct{db.LiveDataBase.GetCollection("live_list")},
	}
}

func (a *LiveApi) RegisterRouter(server *GinServer) {
	// RegisterRouter /******************Start 注册直播间路由*******************/
	server.POST("/live/createLive", a.CreateLive)
	//server.GET("/getAllLive", a.GetAllLive)
	server.GET("/live/fuzzySearchLive", a.FuzzySearchLive)
	server.GET("/live/getRecordList", a.GetRecordList)
	// RegisterRouter /******************End 注册直播间路由*******************/
}

func (a *LiveApi) CreateLive(c *gin.Context) {
	name := c.PostForm("name")
	poster := c.PostForm("poster")
	id := settings.GenNewID()

	filter := bson.D{{"name", name}}

	if !a.LiveColl.CheckDBContains(filter) {
		streamID := settings.ServiceName + "/" + strconv.Itoa(id)
		newLive := settings.Live{
			Name:      name,
			Poster:    poster,
			StartTime: time.Now().Round(time.Second),
			RtmpAddr:  config.Config.RtmpPushPullURL + streamID,
			StreamID:  streamID,
		}

		if a.LiveColl.InsertLive(&newLive) {
			c.JSON(http.StatusOK, newLive)
			fmt.Println(name, poster, streamID)

		} else {
			c.JSON(http.StatusNotAcceptable, "DB Error")
		}

	} else {
		c.JSON(http.StatusConflict, "Live already found")
	}

}

func (a *LiveApi) GetAllLive(c *gin.Context) {
	filter := bson.D{{}}
	sort := bson.D{{"StartTime", 1}}

	results, err := a.LiveColl.FindLive(filter, sort)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, results)
}

func (a *LiveApi) FuzzySearchLive(c *gin.Context) {
	name := c.Query("name")
	filter := bson.D{bson.E{Key: "name",
		Value: bson.M{"$regex": primitive.Regex{Pattern: ".*" + name + ".*", Options: "i"}}}}
	sort := bson.D{{"StartTime", 1}}

	res, err := a.LiveColl.FindLive(filter, sort)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *LiveApi) GetRecordList(c *gin.Context) {
	var files []string
	files, err := a.GetAllFile(config.Config.M7sRecordDir, files)
	fmt.Println(files)
	if err != nil {
		fmt.Println("Error reading the dir: ", err)
		c.JSON(http.StatusInternalServerError, err)

	} else {
		res := make(map[int]string)
		for k, v := range files {
			res[k] = v
		}
		c.JSON(http.StatusOK, res)
	}
}

func (a *LiveApi) GetAllFile(pathname string, s []string) ([]string, error) {
	rd, err := os.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = a.GetAllFile(fullDir, s)
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
