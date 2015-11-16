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

	videoList, err := cl.GetVideoList(json)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Create("output.test")
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString(videoList.title + "\n")
	for _, video := range videoList.videos {
		file.WriteString(video.url + "\n")
		file.WriteString(video.quality + "\n")
		file.WriteString(video.videoType + "\n")
	}
}
