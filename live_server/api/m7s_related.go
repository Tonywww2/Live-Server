package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"live_server/settings"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func PushVideoToStream(c *gin.Context) {
	streamID := c.PostForm("streamID")
	path := c.PostForm("path")

	filterGet := bson.D{{"stream_id", streamID}}
	var result settings.Live
	err := coll.FindOne(context.TODO(), filterGet).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, "Not found in DB")
			return
		}
		panic(err)
	}

	// Create stream on m7s
	params := url.Values{}
	params.Set("streamPath", result.StreamID)
	params.Set("dump", path)
	parseURL, err := url.Parse(settings.CreateStreamURL + strings.Split(path, ".")[1])
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
	StartRecording(streamID, "flv")
	fmt.Println("Start Recording " + streamID)

	filter := bson.D{{"stream_id", streamID}}
	update := bson.D{{"$set", result}}
	res, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	c.String(http.StatusOK, string(b))
	fmt.Println("Update: ", res)

}

func PushStreamToRtmp(c *gin.Context) {
	streamID := c.PostForm("stream_id")
	rtmpAddr := c.PostForm("rtmp_addr")

	filterGet := bson.D{{"streamID", streamID}}
	var result settings.Live
	err := coll.FindOne(context.TODO(), filterGet).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, "Not found in DB")
			return
		}
		panic(err)
	}

	fmt.Println(result)

	if result.IsStreamed {
		// Create stream on m7s
		params := url.Values{}
		parseURL, err := url.Parse(settings.PushURL)
		if err != nil {
			log.Println("err")
		}
		params.Set("target", rtmpAddr)
		params.Set("streamPath", streamID)

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
		c.String(http.StatusNotFound, "Live not streamed yet")

	}

}

func EndStreamAPI(c *gin.Context) {
	streamPath := c.PostForm("streamPath")
	tp := c.PostForm("type")
	if tp == "" {
		tp = "flv"
	}

	StopRecording(streamPath, tp)

	params := url.Values{}
	params.Set("streamPath", streamPath)
	parseURL, err := url.Parse(settings.EndStreamURL)
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

func StartRecording(streamPath string, tp string) {
	// Start record on m7s
	params := url.Values{}
	params.Set("streamPath", streamPath)
	if tp == "" {
		tp = "flv"
	}
	params.Set("type", tp)
	parseURL, err := url.Parse(settings.RecordStartURL)
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

func StopRecording(streamPath string, tp string) {
	// Start record on m7s
	params := url.Values{}
	params.Set("id", streamPath+"/"+tp)
	parseURL, err := url.Parse(settings.RecordStopURL)
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
