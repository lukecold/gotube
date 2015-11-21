package gotube

import (
	"errors"
	"net/url"
	"strconv"
	. "strings"
)

/*
* Get the top k video id from the search result.
* If k is larger than the number of search result, this function would return all the search result
 */
func GetTopKVideoIds(keywords string, k int) ([]string, error) {
	num := 0
	pageNum := 1
	set := make(map[string]bool)
	var idList []string
	for num < k {
		//Get url of search result from #pageNum page
		searchUrl, err := GetSearchUrl(keywords, pageNum)
		if err != nil {
			return idList, err
		}
		//Get list of video id from current page
		idListOfPage, err := GetVideoIdsFromSearch(searchUrl)
		if err != nil {
			return idList, err
		}
		//Add id from id list retrieved in current page to result until we already got top k or out of result
		idIdx := 0
		for num < k && idIdx < len(idListOfPage) {
			_, ok := set[idListOfPage[idIdx]]
			if ok { //We have ran out of search results, it's repeating the last page
				MapToArray(set, &idList)
				return idList, err
			} else { //This id is new
				set[idListOfPage[idIdx]] = true
			}
			idIdx++
			num++
		}
		pageNum++
	}
	MapToArray(set, &idList)
	return idList, nil
}

/*
* Get a search url from the provided keywords
 */
func GetSearchUrl(keywords string, pageNum int) (searchUrl string, err error) {
	//Replace ' ' with '+', like what the YouTube search does
	keywords = Map(
		func(r rune) rune {
			if r == ' ' {
				r = '+'
			}
			return r
		}, keywords)
	//Escape keyword to safely put into url
	keywords = url.QueryEscape(keywords)
	searchUrl = "https://www.youtube.com/results?search_query=" + keywords
	//Make sure page number is valid
	switch {
	case pageNum < 1:
		err = errors.New("invalid page number")
		return
	case pageNum == 1:
		//No action needed
	case pageNum > 1:
		searchUrl += "&page=" + strconv.Itoa(pageNum)
	}
	return
}

/*
* Convert a map[string]bool to string array
 */
func MapToArray(m map[string]bool, a *[]string) {
	for key, _ := range m {
		*a = append(*a, key)
	}
}
