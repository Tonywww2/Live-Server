package main

import (
	"live_server/api"
	"live_server/config"
	"live_server/db"
	"log"
	"os"
)

func main() {
	loadConfig()

	db.InitDB()
	//初始话httpserver
	api.Initialize().RegisterRouters(api.NewLiveApi(), api.NewM7sApi()).Listen()

}

func loadConfig() {
	profile := os.Getenv("Profile")
	log.Printf("current env profile is: %v", profile)
	if profile == "" {
		profile = "Dev"
	}
	//根据环境变量，加载不同环境的配置，比如服务器地址，数据库地址
	switch profile {
	case "Dev":
		config.LoadConfigDev()

	case "Test":
		config.LoadConfigTest()
	}
}
