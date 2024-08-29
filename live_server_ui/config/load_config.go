package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type configuration struct {
	CreateLiveURL      string `json:"create_live_url"`
	UploadUrl          string `json:"upload_url"`
	GetAllLiveURL      string `json:"get_all_live_url"`
	FuzzySearchLiveURL string `json:"fuzzy_search_live_url"`
	ToStreamURL        string `json:"to_stream_url"`
	ToRtmpURL          string `json:"to_rtmp_url"`
	EndStreamUrl       string `json:"end_stream_url"`
	GetRecordsUrl      string `json:"get_records_url"`
	RtmpListUrl        string `json:"rtmp_list_url"`
	RecordsUrl         string `json:"records_url"`
	ImgUrl             string `json:"img_url"`
}

var (
	Config configuration
)

func LoadConfig() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config = configuration{}
	err := decoder.Decode(&Config)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
