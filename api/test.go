package gotube

import (
	"testing"
)

func GetVideoListTesting(t *testing.T) {
	//Testing valid video
	testTitle := "TEST VIDEO"
	cl := Client {VideoRepository: tests}
	vl, err := cl.GetVideoListFromId("C0DPdy98e4c")	//Get test video
	if err != nil {
		t.Fatalf(err)
	}
	if vl.title != testTitle {
		t.Fatalf("Expected title: %v, got: %v", testTitle, vl.title)
	}
	if len(vl.videos) != 5 {
		t.Fatalf("Expected 5 videos, got %v", len(vl.videos))
	}

	//Testing invalid video
	vl, err := cl.GetVideoListFromId("I'm not a valid video id")
	if err == nil {
		t.Fatalf("Expected error for invalid video id")
	}
}

func DownloadTesting(t * testing.T)