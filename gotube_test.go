package gotube_test

import (
	"os"
	"testing"
	. "github.com/KeluDiao/gotube/api"
)

func GetVideoListTesting(t *testing.T) {
	//Testing valid video
	testTitle := "TEST VIDEO"
	cl := Client {VideoRepository: "videos"}
	vl, err := cl.GetVideoListFromId("C0DPdy98e4c")	//Get test video
	if err != nil {
		t.Fatalf(err.Error())
	}
	if vl.Title != testTitle {
		t.Fatalf("Expected title: %v, got: %v", testTitle, vl.Title)
	}
	if len(vl.Videos) != 5 {
		t.Fatalf("Expected 5 videos, got %v", len(vl.Videos))
	}

	//Testing invalid video
	vl, err = cl.GetVideoListFromId("I'm not a valid video id")
	if err == nil {
		t.Fatalf("Expected error for invalid video id")
	}
}

func DownloadTesting(t * testing.T) {
	cl := Client {VideoRepository: "videos"}
	vl, err := cl.GetVideoListFromId("C0DPdy98e4c")	//Get test video
	err = vl.Filter("medium", "video/mp4")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(vl.Videos) != 1 {
		t.Fatalf("Expected 1 videos after filtering, got %v", len(vl.Videos))
	}
	//Download video into ./videos
	err = vl.Download(cl, "", "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	//Remove repository ./videos
	err = os.Remove("videos")
	if err != nil {
		t.Fatalf(err.Error())
	}
}