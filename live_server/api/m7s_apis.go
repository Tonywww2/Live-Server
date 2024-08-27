package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"live_server/db"
	"log"
	"net/http"
	"net/url"
	"strings"

	"live_server/config"
)

type M7sApi struct {
	LiveColl *mongo.Collection
}

// 这里是构造函数
func NewM7sApi() *M7sApi {
	return &M7sApi{
		LiveColl: db.GetCollection("live_list"),
	}
}

// 不同api 处理各自路由
func (a *M7sApi) RegisterRouter(server *GinServer) {
	server.POST("/pushVideoToStream", a.PushVideoToStream)
	server.POST("/pushStreamToRtmp", a.PushStreamToRtmp)
	server.POST("/endStream", a.EndStreamAPI)

	server.StaticFS("/records", http.Dir(config.Config.M7sRecordDir))
}

func (a *M7sApi) PushVideoToStream(c *gin.Context) {
	streamID := c.PostForm("streamID")
	path := c.PostForm("path")

	filterGet := bson.D{{"stream_id", streamID}}
	sort := bson.D{}

	res, err := db.FindLive(a.LiveColl, filterGet, sort)

	if err != nil || len(*res) == 0 {
		c.JSON(http.StatusNotAcceptable, "Invalid Live")
		return
	}

	// Create stream on m7s
	params := url.Values{}
	params.Set("streamPath", (*res)[0].StreamID)
	params.Set("dump", path)
	parseURL, err := url.Parse(config.Config.CreateStreamURL + strings.Split(path, ".")[1])
	if err != nil {
		log.Println("err")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	parseURL.RawQuery = params.Encode()
	urlPathWithParams := parseURL.String()
	resp, err := http.Get(urlPathWithParams)
	if err != nil {
		log.Println("err")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	(*res)[0].IsStreamed = true
	if !a.StartRecording(streamID, "flv") {
		log.Println("Error Starting recording")
		return
	}
	fmt.Println("Start Recording " + streamID)

	filter := bson.D{{"stream_id", streamID}}
	update := bson.D{{"$set", (*res)[0]}}

	if db.UpdateLive(a.LiveColl, filter, update) != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, string(b))

}

func (a *M7sApi) PushStreamToRtmp(c *gin.Context) {
	streamID := c.PostForm("stream_id")
	rtmpAddr := c.PostForm("rtmp_addr")

	filterGet := bson.D{{"stream_id", streamID}}
	sort := bson.D{}
	res, err := db.FindLive(a.LiveColl, filterGet, sort)
	if err != nil || len(*res) == 0 {
		c.JSON(http.StatusNotFound, "Not found in DB")
		log.Println(streamID)
		log.Println(*res)
		log.Println("Not found in DB")
		return
	}

	if (*res)[0].IsStreamed {
		// Create stream on m7s
		params := url.Values{}
		parseURL, err := url.Parse(config.Config.PushURL)
		if err != nil {
			log.Println("err")
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		params.Set("target", rtmpAddr)
		params.Set("streamPath", streamID)

		parseURL.RawQuery = params.Encode()
		urlPathWithParams := parseURL.String()
		resp, err := http.Get(urlPathWithParams)
		if err != nil {
			log.Println("err")
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("err")
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, string(b))

	} else {
		c.String(http.StatusNotFound, "Live not streamed yet")
		log.Println("Live not streamed yet")

	}

}

func (a *M7sApi) EndStreamAPI(c *gin.Context) {
	streamPath := c.PostForm("streamPath")
	tp := c.PostForm("type")
	if tp == "" {
		tp = "flv"
	}

	if !a.StopRecording(streamPath, tp) {
		log.Println("err")
		return
	}

	params := url.Values{}
	params.Set("streamPath", streamPath)
	parseURL, err := url.Parse(config.Config.EndStreamURL)
	if err != nil {
		log.Println("err")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	parseURL.RawQuery = params.Encode()
	urlPathWithParams := parseURL.String()
	fmt.Println(urlPathWithParams)
	resp, err := http.Get(urlPathWithParams)
	if err != nil {
		log.Println("err")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("err")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(string(b))
	fmt.Println("Stream ended: " + streamPath)
	c.String(http.StatusOK, string(b))

}

func (a *M7sApi) StartRecording(streamPath string, tp string) bool {
	// Start record on m7s
	params := url.Values{}
	params.Set("streamPath", streamPath)
	if tp == "" {
		tp = "flv"
	}
	params.Set("type", tp)
	parseURL, err := url.Parse(config.Config.RecordStartURL)
	if err != nil {
		log.Println("err")
		return false
	}

	parseURL.RawQuery = params.Encode()
	urlPathWithParams := parseURL.String()
	fmt.Println(urlPathWithParams)
	resp, err := http.Get(urlPathWithParams)
	if err != nil {
		log.Println("err")
		return false
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("err")
		return false
	}
	fmt.Println(string(b))
	return true
}

func (a *M7sApi) StopRecording(streamPath string, tp string) bool {
	// Stop record on m7s
	params := url.Values{}
	params.Set("id", streamPath+"/"+tp)
	parseURL, err := url.Parse(config.Config.RecordStopURL)
	if err != nil {
		log.Println("err")
		return false
	}

	parseURL.RawQuery = params.Encode()
	urlPathWithParams := parseURL.String()
	fmt.Println(urlPathWithParams)
	resp, err := http.Get(urlPathWithParams)
	if err != nil {
		log.Println("err")
		return false
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err")
		return false
	}
	fmt.Println(string(b))
	return true
}
