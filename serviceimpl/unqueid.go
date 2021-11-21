package serviceimpl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Response struct {
	Uniqueid string
}

func Genereteid() int64 {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://snowflake-vn3k2k7q7q-uc.a.run.app/generateId", nil)
	if err != nil {
		fmt.Print("Error in Generating Unique ID for API call ", err.Error())
		return 0
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("Error in Generating Unique ID In getting res ", err.Error())
		return 0
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("Error in Generating Unique ID parsing response ", err.Error())
		return 0
	}
	var responseObject Response
	json.Unmarshal(bodyBytes, &responseObject)
	id, err := strconv.Atoi(responseObject.Uniqueid)
	if err != nil {
		fmt.Println("Error in Id generation", err.Error())
		return 0
	}
	uniqueid := int64(id)
	return uniqueid
}
