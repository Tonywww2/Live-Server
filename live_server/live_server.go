package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	LiveList        = make(map[string]Live)
	serviceName     = "live"
	rtmpPushPullURL = "rtmp://localhost/"
	createStreamURL = "http://localhost:8080/api/replay/"
	pushURL         = "http://localhost:8080/rtmp/api/push"
	recordStartURL  = "http://localhost:8080/record/api/start"
	recordStopURL   = "http://localhost:8080/record/api/stop"
	endStreamURL    = "http://localhost:8080/api/closestream"
	MongodbUri      = "mongodb://admin:admin@localhost:27017/?retryWrites=true&w=majority"
	coll            *mongo.Collection
)

type Live struct {
	Name       string    `bson:"name"`
	Poster     string    `bson:"poster, omitempty"`
	StartTime  time.Time `bson:"start_time"`
	RtmpAddr   string    `bson:"rtmp_addr"`
	StreamID   string    `bson:"stream_id"`
	IsStreamed bool      `bson:"is_streamed"`
}

func (obj *Live) toString() string {
	return "{Name=" + obj.Name + ", rtmp=" + obj.RtmpAddr + "}"
}

func main() {
	// DB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MongodbUri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll = client.Database("live_list").Collection("live_list")

	// Init Map
	filter := bson.D{{}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	// Unpacks the cursor into a slice
	var result []Live
	if err = cursor.All(context.TODO(), &result); err != nil {
		panic(err)
	}
	for _, v := range result {
		LiveList[v.StreamID] = v
	}
	fmt.Println(LiveList)

	r := gin.Default()
	r.POST("/createLive", createLive)
	r.GET("/getLiveName", getLiveName)
	r.GET("/getAllLive", getAllLive)
	r.POST("/pushVideoToStream", pushVideoToStream)
	r.POST("/pushStreamToRtmp", pushStreamToRtmp)
	r.POST("/endStream", endStreamAPI)

	r.Run(":8082")
}

func createLive(c *gin.Context) {
	//types := c.DefaultPostForm("type", "post")
	name := c.PostForm("name")
	poster := c.PostForm("poster")

	id := int(time.Now().UnixNano())
	id /= 100 / 41
	id %= 1000000000

	flag := true
	for _, v := range LiveList {
		if v.Name == name {
			c.String(http.StatusNotAcceptable, fmt.Sprintf("Already exist"))
			flag = false
		}
	}
	if flag {
		streamID := serviceName + "/" + strconv.Itoa(id)
		LiveList[streamID] = Live{
			name,
			poster,
			time.Now().Round(time.Second),
			rtmpPushPullURL + streamID,
			streamID,
			false,
		}

		c.JSON(http.StatusOK, LiveList[streamID])
		fmt.Println(name, poster, streamID)

		result, err := coll.InsertOne(context.TODO(), LiveList[streamID])
		if err != nil {
			return
		}
		fmt.Println("Inserted ID: %s", result.InsertedID)

	}

}

func getLiveName(c *gin.Context) {
	id := c.DefaultQuery("stream_id", "")
	fmt.Println(id)
	if result, contains := LiveList[id]; contains {
		c.String(http.StatusOK, result.toString())

	} else {
		c.String(http.StatusNotFound, "")

	}

}

func getAllLive(c *gin.Context) {
	c.JSON(http.StatusOK, LiveList)
}

func pushVideoToStream(c *gin.Context) {
	streamID := c.PostForm("streamID")
	path := c.PostForm("path")

	result, ok := LiveList[streamID]

	if ok {
		// Create stream on m7s
		params := url.Values{}
		params.Set("streamPath", result.StreamID)
		params.Set("dump", path)
		parseURL, err := url.Parse(createStreamURL + strings.Split(path, ".")[1])
		if err != nil {
			log.Println("err")
		}

		parseURL.RawQuery = params.Encode()
		urlPathWithParams := parseURL.String()
		fmt.Println(urlPathWithParams)
		resp, err := http.Get(urlPathWithParams)
		if err != nil {
			log.Println("err")
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("err")
		}

		result.IsStreamed = true
		LiveList[streamID] = result
		c.String(http.StatusOK, string(b))

		startRecording(streamID, "flv")
		fmt.Println("Start Recording " + streamID)

		filter := bson.D{{"name", LiveList[streamID].Name}}

		update := bson.D{{"$set", LiveList[streamID]}}

		res, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			fmt.Println("err", err)
			return
		}
		fmt.Println("Update: ", res)

	} else {
		c.String(http.StatusNotFound, "")

	}

}

func pushStreamToRtmp(c *gin.Context) {
	streamPath := c.PostForm("stream_id")
	rtmpAddr := c.PostForm("rtmp_addr")

	result, ok := LiveList[streamPath]

	fmt.Println(result)

	if ok && result.IsStreamed {
		// Create stream on m7s
		params := url.Values{}
		parseURL, err := url.Parse(pushURL)
		if err != nil {
			log.Println("err")
		}
		params.Set("target", rtmpAddr)
		params.Set("streamPath", streamPath)

		parseURL.RawQuery = params.Encode()
		urlPathWithParams := parseURL.String()
		resp, err := http.Get(urlPathWithParams)
		if err != nil {
			log.Println("err")
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("err")
		}

		c.String(http.StatusOK, string(b))

	} else {
		c.String(http.StatusNotFound, "")

	}

}

func startRecording(streamPath string, tp string) {
	// Start record on m7s
	params := url.Values{}
	params.Set("streamPath", streamPath)
	if tp == "" {
		tp = "flv"
	}
	params.Set("type", tp)
	parseURL, err := url.Parse(recordStartURL)
	if err != nil {
		log.Println("err")
	}

	parseURL.RawQuery = params.Encode()
	urlPathWithParams := parseURL.String()
	fmt.Println(urlPathWithParams)
	resp, err := http.Get(urlPathWithParams)
	if err != nil {
		log.Println("err")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err")
	}
	fmt.Println(string(b))
}

func stopRecording(streamPath string, tp string) {
	// Start record on m7s
	params := url.Values{}
	params.Set("id", streamPath+"/"+tp)
	parseURL, err := url.Parse(recordStopURL)
	if err != nil {
		log.Println("err")
	}

	parseURL.RawQuery = params.Encode()
	urlPathWithParams := parseURL.String()
	fmt.Println(urlPathWithParams)
	resp, err := http.Get(urlPathWithParams)
	if err != nil {
		log.Println("err")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err")
	}
	fmt.Println(string(b))

}

func endStreamAPI(c *gin.Context) {
	streamPath := c.PostForm("streamPath")
	tp := c.PostForm("type")
	if tp == "" {
		tp = "flv"
	}

	stopRecording(streamPath, tp)

	params := url.Values{}
	params.Set("streamPath", streamPath)
	parseURL, err := url.Parse(endStreamURL)
	if err != nil {
		log.Println("err")
	}

	parseURL.RawQuery = params.Encode()
	urlPathWithParams := parseURL.String()
	fmt.Println(urlPathWithParams)
	resp, err := http.Get(urlPathWithParams)
	if err != nil {
		log.Println("err")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err")
	}
	fmt.Println(string(b))
	fmt.Println("Stream ended: " + streamPath)

}
