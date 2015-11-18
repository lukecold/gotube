package gotube

import (
	"flag"
	"log"
	"fmt"
)

func main() {
	isDownload := flag.Bool("d", false, "use this flag to download video")
	isRetList := flag.Bool("l", false, "use this flag to retrieve video list")
	url := flag.String("url", "", "video url")
	id := flag.String("id", "", "video id")
	rep := flag.String("-VideoRepository", "", "(optional) repository to store videos (please use relative path)")
	flag.StringVar(rep, "rep", "", "(optional) repository to store videos (please use relative path)")
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
	if *url == "" && *id == "" {
		fmt.Println("Please specify either url and id.")
		invalidCommand = true
	} else if *url != "" && *id != "" {
		fmt.Println("Please don't specify both url and id.")
		invalidCommand = true
	}
	if invalidCommand {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}

	//var url = "https://www.youtube.com/watch?v=6LZM3_wp2ps"
	cl := Client{VideoRepository: *rep}
	var vl VideoList
	var err error
	if *url != "" {
		vl, err = cl.GetVideoListFromUrl(*url)
	} else {
		vl, err = cl.GetVideoListFromId(*id)
	}
	if err != nil {
		log.Fatal(err)
	}
	if *isDownload {
		err = vl.Download(cl, *quality, *extension)
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
