package gotube

import (
	"fmt"
	"os"
	"strconv"
	. "strings"
)

/*
* YouTube video
 */
type Video struct {
	Title     string
	url       string
	quality   string
	extension string
}

/*
* A list of YouTube video with different resolutions
* that shared the same YouTube url.
 */
type VideoList struct {
	Title  string
	Videos []Video
}

/*
* Check if any field is missing.
* A missing filed means a bug found in this program.
 */
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

/*
* Download this video into the repository,
* if repository is not generated, download to current folder.
 */
func (video *Video) Download(rep string) error {
	//Get video from url
	body, err := GetHttpFromUrl(video.url)
	if err != nil {
		return err
	}
	if rep != "" {
		//Make a directory and give every user highest permission
		os.MkdirAll(rep, 0777)
		if !HasSuffix(rep, "/") {
			rep += "/"
		}
	}

	filename := video.Title + video.extension
	//Make sure there is no invalid characters in filename
	filename = Map(
		func(r rune) rune {
			switch r {
			case '/', '\\':
				r = '.'
			case ':', '*', '?', '"', '<', '>', '|':
				r = '-'
			}
			return r
		}, filename)
	filename = rep + filename
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = file.Write(body)
	if err != nil {
		return err
	}
	return nil
}

/*
* Append a video to the video list, video title is assigned here
 */
func (vl *VideoList) Append(v Video) {
	v.Title = vl.Title
	vl.Videos = append(vl.Videos, v)
}

/*
* Download a video from the video list.
* Filter the list first by the given key words,
* then download the first video in the list
 */
func (vl *VideoList) Download(rep string, quality, extension string) (err error) {
	vl.Filter(quality, extension)

	//No matter how many left, pick the first one
	video := vl.Videos[0]
	err = video.Download(rep)
	return err
}

/*
* Filter the video list by given key words.
* The videos don't match are removed from list.
 */
func (vl *VideoList) Filter(quality, extension string) (err error) {
	var matchingVideos []Video
	//Filter by quality
	if quality != "" {
		for _, video := range vl.Videos {
			if video.quality == quality {
				matchingVideos = append(matchingVideos, video)
			}
		}
		vl.Videos = matchingVideos
	}
	matchingVideos = nil
	//Filter by extension
	if extension != "" {
		for _, video := range vl.Videos {
			if video.extension == extension {
				matchingVideos = append(matchingVideos, video)
			}
		}
		vl.Videos = matchingVideos
	}
	if len(vl.Videos) == 0 {
		err = NoMatchingVideoError{_quality: quality, _extension: extension}
		return
	}
	return
}

/*
* Implemented String interface to output VideoList in a delightful format
 */
func (vl VideoList) String() string {
	var videoListStr string
	videoListStr += fmt.Sprintf("video Title: " + vl.Title + "\n")
	videoListStr += fmt.Sprintf("Index\tquality\textension\n")
	for idx, video := range vl.Videos {
		videoListStr += fmt.Sprintf(" %v\t%v\t%v\n",
			strconv.Itoa(idx),
			video.quality,
			video.extension)
	}
	return videoListStr
}
