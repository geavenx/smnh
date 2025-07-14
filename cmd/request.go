package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type CollageRequest struct {
	Rows      string
	Columns   string
	Username  string
	Method    string
	Period    string
	Artist    bool
	Album     bool
	Playcount bool
}

func Request(request CollageRequest) {
	base, _ := url.Parse("https://songstitch.art/collage")

	params := url.Values{}
	params.Add("rows", request.Rows)
	params.Add("columns", request.Columns)
	params.Add("username", request.Username)
	params.Add("method", string(request.Method))
	params.Add("period", string(request.Period))
	params.Add("artist", fmt.Sprintf("%t", request.Artist))
	params.Add("album", fmt.Sprintf("%t", request.Album))
	params.Add("playcount", fmt.Sprintf("%t", request.Playcount))
	base.RawQuery = params.Encode()

	resp, err := http.Get(base.String())
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("Bad status: " + resp.Status)
	}

	out, err := os.Create("smnh.jpeg")
	if err != nil {
		panic(err)
	}

	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		panic(err)
	}

	fmt.Println("image saved to smnh.jpeg")
}
