package gotube

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	. "strings"
)

/*
* YouTube client.
 */
type Client struct {
	VideoRepository string
}

/*
* Get a video list from given id.
 */
func (cl *Client) GetVideoListFromId(id string) (VideoList, error) {
	url := "https://www.youtube.com/watch?v=" + id
	return cl.GetVideoListFromUrl(url)
}

/*
* Get a video list from given url.
 */
func (cl *Client) GetVideoListFromUrl(url string) (vl VideoList, err error) {
	//Get webpage content from url
	body, err := cl.GetHttpFromUrl(url)
	if err != nil {
		return
	}
	//Extract json data from webpage content
	jsonData, err := cl.GetJsonFromHttp(body)
	if err != nil {
		return
	}
	//Fetch video list according to json data
	vl, err = cl.GetVideoListFromJson(jsonData)
	if err != nil {
		return
	}
	return
}

/*
* Initialize a GET request, and get the http code of the webpage.
 */
func (cl *Client) GetHttpFromUrl(url string) ([]byte, error) {
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
* Get json data from http code.
 */
func (*Client) GetJsonFromHttp(httpData []byte) (map[string]interface{}, error) {
	//Find out if this page is age-restricted
	if bytes.Index(httpData, []byte("og:restrictions:age")) != -1 {
		return nil, AgeRestrictedError{}
	}
	//Find begining of json data
	var jsonBeg = "ytplayer.config = {"
	beg := bytes.Index(httpData, []byte(jsonBeg))
	if beg == -1 { //pattern not found
		return nil, PatternNotFoundError{_pattern: jsonBeg}
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
* Get video list from json data retrieved from http code.
 */
func (*Client) GetVideoListFromJson(jsonData map[string]interface{}) (vl VideoList, err error) {
	args := jsonData["args"].(map[string]interface{})
	vl.Title = args["title"].(string)
	encodedStreamMap := args["url_encoded_fmt_stream_map"].(string)
	//Videos are seperated by ","
	videoListStr := Split(encodedStreamMap, ",")
	for _, videoStr := range videoListStr {
		//Parameters of a video are seperated by "&"
		videoParams := Split(videoStr, "&")
		var video Video
		for _, param := range videoParams {
			/*Unescape the url encoding characters.
			Only do it after seperation because
			there are "," and "&" escaped in url*/
			param, err = url.QueryUnescape(param)
			if err != nil {
				return
			}
			switch {
			case HasPrefix(param, "quality"):
				video.quality = param[8:]
			case HasPrefix(param, "type"):
				//type and codecs are seperated by ";"
				video.extension = Split(param, ";")[0][5:]
			case HasPrefix(param, "url"):
				video.url = param[4:]
			}
		}
		missingFields := video.FindMissingFields()
		if missingFields != nil {
			err = MissingFieldsError{_fields: missingFields}
			return
		}
		vl.Append(video)
	}
	return
}
