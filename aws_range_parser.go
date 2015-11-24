package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"encoding/json"
)

type Ranges struct {
	CreateDate string `json:"createDate"`
	Prefixes   []struct {
		IPPrefix string `json:"ip_prefix"`
		Region   string `json:"region"`
		Service  string `json:"service"`
	} `json:"prefixes"`
	SyncToken string `json:"syncToken"`
}

func main() {
	response, err := http.Get("https://ip-ranges.amazonaws.com/ip-ranges.json")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		jsonSrc := []byte(string(contents))
		var myJson Ranges
		json.Unmarshal(jsonSrc, &myJson)
		if myJson.Prefixes {
			fmt.Println("Example")
		}
//		fmt.Printf("%s\n", string(contents))
	}
}
