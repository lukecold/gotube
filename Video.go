package main

/*
YouTube video
*/
type Video struct {
	url       string
	quality   string
	extension string
}

type VideoList struct {
	title  string
	videos []Video
}

func (vl *VideoList) Append(v Video) {
	vl.videos = append(vl.videos, v)
}
