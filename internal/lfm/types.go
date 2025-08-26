package lfm

import "encoding/xml"

type Period string
type Method string

const (
	PeriodOverall Period = "overall"
	Period7Day    Period = "7day"
	Period1Month  Period = "1month"
	Period3Month  Period = "3month"
	Period6Month  Period = "6month"
	Period12Month Period = "12month"

	TopAlbums  Method = "user.getTopAlbums"
	TopArtists Method = "user.getTopArtists"
	TopTracks  Method = "user.getTopTracks"
)

type Image struct {
	XMLName xml.Name `xml:"image"`
	Size    string   `xml:"size,attr"`
	Url     string   `xml:",chardata"`
}

type ImageResult struct {
	Filename string
	Err      error
}
