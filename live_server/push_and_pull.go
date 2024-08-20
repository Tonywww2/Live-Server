package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func vtsGet(wUrl string, file string, streamPath string) (result string) {
	params := url.Values{}
	params.Set("dump", file)
	params.Set("streamPath", streamPath)
	parseURL, err := url.Parse(wUrl + strings.Split(file, ".")[1])
	if err != nil {
		log.Println("err")
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	parseURL.RawQuery = params.Encode()
	urlPathWithParams := parseURL.String()
	resp, err := http.Get(urlPathWithParams)
	if err != nil {
		log.Println("err")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err")
	}

	result = string(b)

	return result

}

func strGet(wUrl string, target string, streamPath string) (result string) {
	params := url.Values{}
	parseURL, err := url.Parse(wUrl)
	if err != nil {
		log.Println("err")
	}
	params.Set("target", target)
	params.Set("streamPath", streamPath)
	//如果参数中有中文参数,这个方法会进行URLEncode
	parseURL.RawQuery = params.Encode()
	urlPathWithParams := parseURL.String()
	resp, err := http.Get(urlPathWithParams)
	if err != nil {
		log.Println("err")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err")
	}

	result = string(b)

	return result

}

func formatAddr(addr string) string {
	if addr[0] == '"' {
		addr = addr[1:]
		addr = addr[:len(addr)-1]
	}
	return addr
}

var (
	vtsUrl = "http://localhost:8080/api/replay/"
	strUrl = "http://localhost:8080/rtmp/api/push"
)

func main() {
	var fileDir string
	var rtmpAddr string
	var serviceName = "vtos"

	fmt.Println("Enter video path: ")
	n, err := fmt.Scan(&fileDir)

	if err != nil || n > 1 {
		fmt.Println("Wrong Format")
	}

	fmt.Println("Enter rtmp address: ")
	n, err = fmt.Scan(&rtmpAddr)

	if err != nil || n > 1 {
		fmt.Println("Wrong Format")
	}

	fileDir = formatAddr(fileDir)
	rtmpAddr = formatAddr(rtmpAddr)

	id := int(time.Now().UnixNano())
	id /= 100 / 41
	id %= 1000000000

	streamID := serviceName + "/" + strconv.Itoa(id)
	fmt.Println("Stream ID:", streamID)
	fmt.Println(vtsGet(vtsUrl, fileDir, streamID))
	fmt.Println(strGet(strUrl, rtmpAddr, streamID))

}
