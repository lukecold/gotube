package gotube

import (
	"io/ioutil"
	"net/http"
)

/*
* Initialize a GET request, and get the http code of the webpage.
 */
func GetHttpFromUrl(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

/*
* Get a video list from given id.
 */
func GetVideoListFromId(id string) (VideoList, error) {
	url := "https://www.youtube.com/watch?v=" + id
	return GetVideoListFromUrl(url)
}

/*
* Get a video list from given url.
 */
func GetVideoListFromUrl(url string) (vl VideoList, err error) {
	//Get webpage content from url
	body, err := GetHttpFromUrl(url)
	if err != nil {
		return
	}
	//Extract json data from webpage content
	jsonData, err := GetJsonFromHttp(body)
	if err != nil {
		return
	}
	//Fetch video list according to json data
	vl, err = GetVideoListFromJson(jsonData)
	if err != nil {
		return
	}
	return
}