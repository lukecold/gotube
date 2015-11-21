package main

import (
	"flag"
	"fmt"
	. "github.com/KeluDiao/gotube/api"
	"log"
	"runtime"
	"sync"
)

func main() {
	isDownload := flag.Bool("d", false, "use this flag to download video")
	isRetList := flag.Bool("l", false, "use this flag to retrieve video list")
	url := flag.String("url", "", "video url")
	id := flag.String("id", "", "video id")
	search := flag.String("-search", "", "search by key words (specify top k with command -k)")
	flag.StringVar(search, "s", "", "search by key words (specify top k with command -k)")
	k := flag.Int("k", 1, "return top k results, only valid with key word searching")
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
		ids, err := GetTopKVideoIds(*search, *k)
		if err != nil {
			log.Fatal(err)
		}
		if *isRetList {
			fmt.Printf("The top %v results for key words \"%v\" are:\n\n", *k, *search)
		}
		//Waiting group is used to prevent main thread ending before child threads end
		wg := new(sync.WaitGroup)
		wg.Add(len(ids))
		//Channel is used to control the maximum threads
		end := make(chan bool, MaxParallelism())
		for _, vid := range ids {
			vl, err = GetVideoListFromId(vid)
			if err != nil {
				log.Fatal(err)
			}
			go Exec(vl, *isDownload, *isRetList, *rep, *quality, *extension, wg, end)
			end <- true
		}
		wg.Wait()
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	//dummy variables
	wg := new(sync.WaitGroup)
	wg.Add(1)
	var end chan bool
	Exec(vl, *isDownload, *isRetList, *rep, *quality, *extension, wg, end)
}

/*
* Choose either downloading or retrieving video list
*/
func Exec(vl VideoList, isDownload, isRetList bool, rep, quality, extension string, wg *sync.WaitGroup, end chan bool) {
	//Set up synchronization function
	defer func() {
		<- end
		wg.Done()
	}

	if isDownload {
		fmt.Printf("Downloading %v...\n", vl.Title)
		err := vl.Download(rep, quality, extension)
		if err != nil {
			log.Fatal(err)
		}
	} else if isRetList {
		fmt.Printf("Videos under the name of %v:\n", vl.Title)
		err := vl.Filter(quality, extension)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(vl)
	}
}

/*
* Find out the maximum go routines allowed
*/
func MaxParallelism() int {
    maxProcs := runtime.GOMAXPROCS(0)
    numCPU := runtime.NumCPU()
    if maxProcs < numCPU {
        return maxProcs
    }
    return numCPU
}