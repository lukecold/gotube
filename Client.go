package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

/*
YouTube client
You need to log in to view age-restricted videos
*/
type Client struct {
	userName string
	passWord string
}

/*
YouTube video
*/
type Video struct {
	title string
}

/*
Request http code from url
*/
func (*Client) RequestUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

/*
Get json data
*/
func (*Client) GetJson(httpData []byte) (map[string]interface{}, error) {
	//Find begining of json data
	var jsonBeg = "ytplayer.config = {"
	beg := bytes.Index(httpData, []byte(jsonBeg))
	if beg == -1 { //pattern not found
		return nil, PatternNotFoundError{pattern: jsonBeg}
	}
	beg += len(jsonBeg) //len(jsonBeg) returns the number of bytes in jsonBeg

	//Find offset of json data
	var unmatchedBrackets = 1
	var offset = 0
	for unmatchedBrackets > 0 {
		nextRight := bytes.Index(httpData[beg+offset:], []byte("}"))
		if nextRight == -1 {
			return nil, UnmatchedBracketsError{}
		}
		unmatchedBrackets -= 1
		unmatchedBrackets += bytes.Count(httpData[beg+offset:beg+offset+nextRight], []byte("{"))
		offset += nextRight + 1
	}

	//Load json data
	var f interface{}
	err := json.Unmarshal(httpData[beg-1:beg+offset], &f)
	if err != nil {
		return nil, err
	}
	return f.(map[string]interface{}), nil
}

/*
Get video from json data
*/
func (*Client) GetVideo(jsonData map[string]interface{}) (video Video, err error) {
	args := jsonData["args"].(map[string]interface{})
	video.title = args["title"].(string)
	

	return
}
