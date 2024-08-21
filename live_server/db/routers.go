package db

import "go.mongodb.org/mongo-driver/mongo"

var (
	LiveColl *mongo.Collection
)

func InitRouters() {
	LiveColl = GetCollection("live_list")

}
