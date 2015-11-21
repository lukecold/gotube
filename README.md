# gotube ![build status](https://travis-ci.org/KeluDiao/gotube.svg?branch=master)
Gotube is a YouTube downloader using go language.
go language is a new light-weight language developed by Google, 
it provides not only many powerful libraries, but also a simple multi-threading syntax.

This tool is an easy way to download any non-age-restricted videos in YouTube. 
You can also perform batch downloading by keywords via search function.
Gotube will generate a number of go-routines (no more than the number of your CPU cores) to download multiple videos simultaneously. 

##Installation
- Install go from [https://golang.org/](https://golang.org/)
- Set up go environment as in [https://golang.org/doc/install](https://golang.org/doc/install)
  - (For Un*x users) 
  	- ```export GOPATH=$HOME/path/to/workspace/```
  	- ```export PATH=$PATH:$GOPATH/bin```
  - (For Windows users) 
  	- In environment variables add GOPATH=path/to/workspace/ 
  	- Append path/to/workspace/bin to PATH
- For command-line usage
  - Type the following in command/terminal
    - ```go get github.com/KeluDiao/gotube```
    - ```gotube -h```
- For library usage
  - You don't need to do anything

##Command-line usage
You can check the video list from a url or video id:
```
gotube -l -url https://www.youtube.com/watch?v=C0DPdy98e4c
```
Filter the video list by specifying video quality:
```
gotube -l -id C0DPdy98e4c -q medium
```
Download the video with default resolution:
```
gotube -d -url https://www.youtube.com/watch?v=C0DPdy98e4c
```
Download the video with specified resolution:
```
gotube -d -id C0DPdy98e4c -q medium -ext video/mp4
```
Specify a video repository using relative path:
```
gotube -d -id C0DPdy98e4c -q medium -ext video/mp4 -rep ./videos
```
Specify a video repository using absolute path (If you are a Windows user, I highly recommend you using absolute path, or the videos may be downloaded into your users folder):
```
gotube -d -id C0DPdy98e4c -q medium -ext video/mp4 -rep /Users/yourusername/Documents/videos
```
Try search by keywords and see what would return:
```
gotube -l -s "curry highlights"
```
Try get more results by querying top k results explicitly (If k is larger than the number of all videos return, the program would return all the videos):
```
gotube -l -s "curry highlights" -k 5
```
Download all of them:
```
gotube -d -s "curry highlights" -k 5 -rep /Users/yourusername/Documents/videos
```

#Library usage
```go
package main

import (
  . "github.com/KeluDiao/gotube/api"
  "log"
)

func main() {
	idList := [...]string{"shLTrG_noKo", "Ojv7tKpzkyM", "GahnMbhmt7g"}
	rep := "Curry_highlights"
	for _, id := range idList {
		vl, err := GetVideoListFromId(id)
		if err != nil {
			log.Fatal(err)
		}
		err = vl.Download(rep, "", "video/mp4")
		if err != nil {
			log.Fatal(err)
		}
	}
}
```

This program is still under-developing. More interesting functionalities will be added into it soon! 
