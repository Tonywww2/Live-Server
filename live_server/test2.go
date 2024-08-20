package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	response, err := http.Get("http://localhost:8082/getAllLive")
	if err != nil || response.StatusCode != 200 {
		fmt.Println("Error")
	}
	var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
	}

}
