package main

import (
	"net/http"
	"io/ioutil"
)

type Client struct {
	userName string
	passWord string
}

func (*Client) Request(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if(err != nil) {
		return "", err
	}
	return string(body), nil
}