package gotube

import (
	"fmt"
	"os"
	"strconv"
	. "strings"
)

/*
YouTube video
*/
type Video struct {
	Title     string
	url       string
	quality   string
	extension string
}

type VideoList struct {
	Title  string
	videos []Video
}

func (video *Video) FindMissingFields() (missingFields []string) {
	if video.quality == "" {
		missingFields = append(missingFields, "quality")
	}
	if video.extension == "" {
		missingFields = append(missingFields, "video type")
	}
	if video.url == "" {
		missingFields = append(missingFields, "url")
	}
	return
}

func (video *Video) Download(cl Client) error {
	//Get video from url
	body, err := cl.GetHttpFromUrl(video.url)
	if err != nil {
		return err
	}
	var pathname string
	if cl.VideoRepository != "" {
		//Make a directory and give every user highest permission
		os.MkdirAll(cl.VideoRepository, 0777)
		pathname = cl.VideoRepository
		if !HasSuffix(pathname, "/") {
			pathname += "/"
		}
	}

	filename := video.Title + video.extension
	filename = Map(
		func(r rune) rune {
			if r == '/' {
				r = '.'
			}
			return r
		}, filename)
	filename = pathname + filename
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	file.Write(body)
	return nil
}

func (vl *VideoList) Append(v Video) {
	v.Title = vl.Title
	vl.videos = append(vl.videos, v)
}

func (vl *VideoList) Download(cl Client, quality, extension string) (err error) {
	vl.Filter(quality, extension)

	//No matter how many left, pick the first one
	video := vl.videos[0]
	err = video.Download(cl)
	return err
}

func (vl *VideoList) Filter(quality, extension string) (err error) {
	var matchingVideos []Video
	if quality != "" {
		for _, video := range vl.videos {
			if video.quality == quality {
				matchingVideos = append(matchingVideos, video)
			}
		}
		vl.videos = matchingVideos
	}
	matchingVideos = nil
	if extension != "" {
		for _, video := range vl.videos {
			if video.extension == extension {
				matchingVideos = append(matchingVideos, video)
			}
		}
		vl.videos = matchingVideos
	}
	if len(vl.videos) == 0 {
		err = NoMatchingVideoError{_quality: quality, _extension: extension}
		return
	}
	return
}

func (vl VideoList) String() string {
	var videoListStr string
	videoListStr += fmt.Sprintf("video Title: " + vl.Title + "\n")
	videoListStr += fmt.Sprintf("Index\tquality\textension\n")
	for idx, video := range vl.videos {
		videoListStr += fmt.Sprintf(" %v\t%v\t%v\n", 
			strconv.Itoa(idx),
			video.quality,
			video.extension)
	}
	return videoListStr
}