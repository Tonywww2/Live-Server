package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type configuration struct {
	RtmpPushPullURL string `json:"rtmp_push_pull_url"`
	CreateStreamURL string `json:"create_stream_url"`
	PushURL         string `json:"push_url"`
	RecordStartURL  string `json:"record_start_url"`
	RecordStopURL   string `json:"record_stop_url"`
	EndStreamURL    string `json:"end_stream_url"`
	MongodbUri      string `json:"mongodb_uri"`
	Dbname          string `json:"dbname"`
	M7sRecordDir    string `json:"m7s_record_dir"`
	Port            int    `json:"port"`
}

var (
	Config configuration
)

func LoadConfigDev() {
	file, _ := os.Open("config_dev.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config = configuration{}
	err := decoder.Decode(&Config)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func LoadConfigTest() {
	file, _ := os.Open("config_test.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config = configuration{}
	err := decoder.Decode(&Config)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
