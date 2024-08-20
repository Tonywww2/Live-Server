package settings

import (
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	ServiceName     = "live"
	RtmpPushPullURL = "rtmp://localhost/"
	CreateStreamURL = "http://localhost:8080/api/replay/"
	PushURL         = "http://localhost:8080/rtmp/api/push"
	RecordStartURL  = "http://localhost:8080/record/api/start"
	RecordStopURL   = "http://localhost:8080/record/api/stop"
	EndStreamURL    = "http://localhost:8080/api/closestream"
	MongodbUri      = "mongodb://admin:admin@localhost:27017/?retryWrites=true&w=majority"
)

var (
	LiveList = make(map[string]Live)
	Coll     *mongo.Collection
)

type Live struct {
	Name       string    `bson:"name"`
	Poster     string    `bson:"poster, omitempty"`
	StartTime  time.Time `bson:"start_time"`
	RtmpAddr   string    `bson:"rtmp_addr"`
	StreamID   string    `bson:"stream_id"`
	IsStreamed bool      `bson:"is_streamed"`
}

func (obj *Live) ToString() string {
	return "{Name=" + obj.Name + ", rtmp=" + obj.RtmpAddr + "}"
}

func GenNewID() int {
	id := int(time.Now().UnixNano())
	id /= 100 / 41
	id %= 1000000000
	return id

}
