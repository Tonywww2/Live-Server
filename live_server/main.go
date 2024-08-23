package main

import (
	"github.com/gin-gonic/gin"
	"live_server/api"
	"live_server/config"
	"live_server/db"
)

func main() {
	config.LoadConfig()

	db.InitDB()
	db.InitRouters()

	r := gin.Default()
	api.OpenAPIs(r)

	err := r.Run(":8082")
	if err != nil {
		return
	}
}
