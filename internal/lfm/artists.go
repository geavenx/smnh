package lfm

import "encoding/xml"

type Artist struct {
	XMLName xml.Name `xml:"artist"`
	Name    string   `xml:"name"`
	MBID    string   `xml:"mbid"`
	Url     string   `xml:"url"`
}
