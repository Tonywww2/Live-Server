package main

import (
	"github.com/gin-gonic/gin"
	"live_server/api"
)

func main() {
	r := gin.Default()
	r.POST("/createLive", api.CreateLive)
	r.GET("/getAllLive", api.GetAllLive)
	r.POST("/pushVideoToStream", api.PushVideoToStream)
	r.POST("/pushStreamToRtmp", api.PushStreamToRtmp)
	r.POST("/endStream", api.EndStreamAPI)

	r.Run(":8082")
}
