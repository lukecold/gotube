package main

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
)

func main() {
	var url string
	fmt.Scan(&url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(string(body))
}