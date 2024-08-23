package db

import (
	"context"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"

	"live_server/config"
)

var (
	MongoClient   *mongo.Client
	MongoDatabase string
)

func InitDB() {
	mongoUri := config.Config.MongodbUri
	if mongoUri == "" {
		log.Fatal("MongoDB URI is not provided in the configuration")
	}
	dbName := config.Config.Dbname
	if dbName == "" {
		log.Fatal("MongoDB database is not provided in the configuration")
	}
	clientOptions := options.Client().ApplyURI(mongoUri)

	// 打开日志调试，实际情况根据生产，测试环境打开
	clientOptions.SetMonitor(&event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Printf("MongoDB Command Started: %s %v", evt.CommandName, evt.Command)
		},
		Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
			log.Printf("MongoDB Command Succeeded: %s %v", evt.CommandName, evt.Reply)
		},
		Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
			log.Printf("MongoDB Command Failed: %s %v", evt.CommandName, evt.Failure)
		},
	})
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}
	log.Println("Connected to MongoDB!")
	MongoClient = client
	MongoDatabase = dbName
}

// 获取不同的集合名
func GetCollection(collection string) *mongo.Collection {
	return MongoClient.Database(MongoDatabase).Collection(collection)
}
