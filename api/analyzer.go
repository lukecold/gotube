package gotube

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	. "strings"
)

/*
* Get json data from http code.
 */
func GetJsonFromHttp(httpData []byte) (map[string]interface{}, error) {
	//Find out if this page is age-restricted
	if bytes.Index(httpData, []byte("og:restrictions:age")) != -1 {
		return nil, errors.New("this page is age-restricted")
	}
	//Find begining of json data
	jsonBeg := "ytplayer.config = {"
	beg := bytes.Index(httpData, []byte(jsonBeg))
	if beg == -1 { //pattern not found
		return nil, PatternNotFoundError{_pattern: jsonBeg}
	}
	beg += len(jsonBeg) - 1 //len(jsonBeg) returns the number of bytes in jsonBeg

	//Load json data
	var f interface{}
	err := json.NewDecoder(bytes.NewReader(httpData[beg:])).Decode(&f)
	if err != nil {
		return nil, err
	}
	return f.(map[string]interface{}), nil
}

/*
* Get video list from json data retrieved from http code.
 */
func GetVideoListFromJson(jsonData map[string]interface{}) (vl VideoList, err error) {
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

/*
* Parse the http data of the page get from url and retrieve the id list
 */
func GetVideoIdsFromSearch(searchUrl string) (idList []string, err error) {
	//Get the http code of the page get from url
	body, err := GetHttpFromUrl(searchUrl)
	if err != nil {
		return
	}
	//Retrive id list
	idBeg := []byte("class=\"yt-lockup yt-lockup-tile yt-lockup-video vve-check clearfix\" data-context-item-id=\"")
	beg := 0
	for {
		//Find the index of begin pattern
		offset := bytes.Index(body[beg:], idBeg)
		if offset < 0 {
			return
		}
		beg += offset + len(idBeg)
		//Find the index of closing parenthesis
		offset = bytes.Index(body[beg:], []byte("\""))
		if offset < 0 {
			err = errors.New("unmatched parenthesis")
			return
		}
		end := beg + offset
		idList = append(idList, string(body[beg:end]))
	}
	return
}
