package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	var url string
	fmt.Scan(&url)
	var cl Client
	body, err := cl.RequestUrl(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	json, err := cl.GetJson(body)
	if err != nil {
		fmt.Println(err)
		return
	}

	video, err := cl.GetVideo(json)

	file, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	file.Write(video.title)
}
