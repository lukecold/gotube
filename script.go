package main

import (
	"flag"
	"fmt"
	. "github.com/KeluDiao/gotube/api"
	"log"
)

func main() {
	isDownload := flag.Bool("d", false, "use this flag to download video")
	isRetList := flag.Bool("l", false, "use this flag to retrieve video list")
	url := flag.String("url", "", "video url")
	id := flag.String("id", "", "video id")
	search := flag.String("-search", "", "search by key words")
	flag.StringVar(search, "s", "", "search by key words")
	rep := flag.String("-videorepository", "", "(optional) repository to store videos")
	flag.StringVar(rep, "rep", "", "(optional) repository to store videos")
	quality := flag.String("-quality", "", "(optional) video quality. e.g. medium")
	flag.StringVar(quality, "q", "", "(optional) video quality. e.g. medium")
	extension := flag.String("-extension", "", "(optional) video extension. e.g. video/mp4, video/flv, video/webm")
	flag.StringVar(extension, "ext", "", "(optional) video extension. e.g. mp4, flv")
	isHelp := flag.Bool("h", false, "help")

	flag.Parse()

	if *isHelp {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}
	invalidCommand := false
	if *isDownload && *isRetList {
		fmt.Println("You can only either download or retrieve video list.")
		invalidCommand = true
	} else if !*isDownload && !*isRetList {
		fmt.Println("Please choose if you want to download or retrieve video list.")
		invalidCommand = true
	}
	//Find out how many sources are specified
	sourceNum := 0
	if *url != "" {
		sourceNum++
	}
	if *id != "" {
		sourceNum++
	}
	if *search != "" {
		sourceNum++
	}
	if sourceNum == 0 {
		fmt.Println("Please specify one of url, id, and key word(s).")
		invalidCommand = true
	} else if sourceNum > 1 {
		fmt.Println("Please don't specify more than one of url, id, and key word(s).")
		invalidCommand = true
	}
	if invalidCommand {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}

	//Get the video list
	var vl VideoList
	var err error
	if *url != "" {
		vl, err = GetVideoListFromUrl(*url)
	} else if *id != "" {
		vl, err = GetVideoListFromId(*id)
	} else {
		ids, err := GetTopKVideoIds(*search, 1)
		if err != nil {
			log.Fatal(err)
		}
		vl, err = GetVideoListFromId(ids[0])
	}
	if err != nil {
		log.Fatal(err)
	}
	//Choose either downloading or retrieving video list
	if *isDownload {
		err = vl.Download(*rep, *quality, *extension)
		if err != nil {
			log.Fatal(err)
		}
	} else if *isRetList {
		err = vl.Filter(*quality, *extension)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(vl)
	}
}
