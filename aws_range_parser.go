package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

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
		fmt.Printf("%s\n", string(contents))
	}
}
