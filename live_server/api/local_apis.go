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
		db.NewLiveColl(),
	}
}

// LiveForm 用于绑定请求参数
type LiveForm struct {
	Name   string `form:"name" binding:"required"`
	Poster string `form:"poster" binding:"required"`
}

func (a *LiveApi) RegisterRouter(server *GinServer) {
	// RegisterRouter /******************Start 注册直播间路由*******************/
	server.POST("/live/createLive", a.CreateLive)
	//server.GET("/getAllLive", a.GetAllLive)
	server.GET("/live/fuzzySearchLive", a.FuzzySearchLive)
	server.GET("/live/getRecordList", a.GetRecordList)
	server.POST("/live/UploadFile", a.UploadFile)
	server.StaticFS("/live_posters", http.Dir("./uploads/live_posters/"))
	// RegisterRouter /******************End 注册直播间路由*******************/
}

func (a *LiveApi) CreateLive(c *gin.Context) {

	var liveForm LiveForm
	if err := c.ShouldBind(&liveForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := settings.GenNewID()

	filter := bson.D{{"name", liveForm.Name}}
	if a.LiveColl.CheckDBContains(filter) {
		c.JSON(http.StatusConflict, "Live already found")
		return
	}

	streamID := settings.ServiceName + "/" + strconv.Itoa(id)
	newLive := settings.Live{
		Name:      liveForm.Name,
		Poster:    liveForm.Poster,
		StartTime: time.Now().Round(time.Second),
		RtmpAddr:  config.Config.RtmpPushPullURL + streamID,
		StreamID:  streamID,
	}

	if a.LiveColl.InsertLive(&newLive) {
		c.JSON(http.StatusOK, newLive)
		fmt.Printf("新直播流创建成功: %s, %s, %s\n", liveForm.Name, liveForm.Poster, streamID)
	} else {
		c.JSON(http.StatusInternalServerError, "数据库错误")
	}

}

//// 数据库服务流量内存，自己的业务服务器内存也吃不消，前端也处理不过来，采用分页查询
//func (a *LiveApi) GetAllLive(c *gin.Context) {
//	filter := bson.D{{}}
//	sort := bson.D{{"StartTime", 1}}
//
//	results, err := a.LiveColl.FindLive(filter, sort)
//
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, results)
//}

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

func (a *LiveApi) UploadFile(c *gin.Context) {
	// 从表单中获取文件
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("获取文件错误: %s", err.Error()))
		return
	}

	savePath := "./uploads/live_posters/" + file.Filename
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("保存文件错误: %s", err.Error()))
		return
	}

	// 返回上传成功信息
	c.String(http.StatusOK, fmt.Sprintf(savePath))
}

// todo 这个接口是有问题的，看着只是读取给定目录的全部文件供客户端选择
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
