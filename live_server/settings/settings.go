package settings

import (
	"time"
)

const (
	ServiceName = "live"
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
