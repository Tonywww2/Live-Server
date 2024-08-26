package api

import (
	"github.com/gin-gonic/gin"
	"live_server/config"
	"net/http"
)

func OpenAPIs(r *gin.Engine) {
	r.POST("/createLive", CreateLive)
	r.GET("/getAllLive", GetAllLive)
	r.GET("/fuzzySearchLive", fuzzySearchLive)
	r.POST("/pushVideoToStream", PushVideoToStream)
	r.POST("/pushStreamToRtmp", PushStreamToRtmp)
	r.POST("/endStream", EndStreamAPI)
	r.GET("/getRecordList", GetRecordList)
	r.StaticFS("/records", http.Dir(config.Config.M7sRecordDir))
}
