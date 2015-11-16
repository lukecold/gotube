package main

import (
	"fmt"
)

func main() {
	var url string
	fmt.Scan(&url)
	var cl Client
	body, err := cl.Request(url)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(body)
	}
}