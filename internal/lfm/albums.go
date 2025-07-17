package lfm

import (
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Period string

const (
	PeriodOverall Period = "overall"
	Period7Day    Period = "7day"
	Period1Month  Period = "1month"
	Period3Month  Period = "3month"
	Period6Month  Period = "6month"
	Period12Month Period = "12month"
)

type TopAlbumRequest struct {
	Username string
	Period   Period
	Limit    int
}

type TopAlbumResponse struct {
	XMLName   xml.Name `xml:"lfm"`
	Status    string   `xml:"status,attr"`
	TopAlbums TopAlbum `xml:"topalbums"`
}

type Album struct {
	XMLName   xml.Name `xml:"album"`
	Rank      string   `xml:"rank,attr"`
	Name      string   `xml:"name"`
	PlayCount int      `xml:"playcount"`
	MBID      string   `xml:"mbid"`
	Url       string   `xml:"url"`
	Artist    Artist   `xml:"artist"`
	Images    []Image  `xml:"image"`
}

type TopAlbum struct {
	XMLName xml.Name `xml:"topalbums"`
	User    string   `xml:"user,attr"`
	Type    string   `xml:"type,attr"`
	Albums  []Album  `xml:"album"`
}

type Image struct {
	XMLName xml.Name `xml:"image"`
	Size    string   `xml:"size,attr"`
	Url     string   `xml:",chardata"`
}

type ImageResult struct {
	Filename string
	Err      error
}

func FetchTopAlbums(req TopAlbumRequest) []Album {
	if err := godotenv.Load(); err != nil {
		slog.Error("error loading .env", "error", err)
	}

	baseUrl := os.Getenv("LAST_FM_BASE_URL")
	apiKey := os.Getenv("LAST_FM_API_KEY")

	apiUrl, err := url.Parse(baseUrl)
	if err != nil {
		slog.Error("error parsing Lfm api url", "error", err)
	}

	params := map[string]string{
		"method":  "user.getTopAlbums",
		"user":    req.Username,
		"api_key": apiKey,
		"period":  string(req.Period),
		"limit":   strconv.Itoa(req.Limit),
	}

	q := url.Values{}
	for k, v := range params {
		q.Set(k, v)
	}

	apiUrl.RawQuery = q.Encode()

	slog.Debug("Sending request", "method", "user.getTopAlbums", "url", apiUrl.String())
	resp, err := http.Get(apiUrl.String())
	if err != nil {
		slog.Error("Failed to execute request", "error", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body", "error", err)
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("API returned no OK status", "status", resp.StatusCode, "body", string(body))
	}

	slog.Debug("response OK", "method", params["method"], "body", string(body))

	var topAlbumsResp TopAlbumResponse
	if err = xml.Unmarshal(body, &topAlbumsResp); err != nil {
		slog.Error("Failed to parse XML", "error", err)
	}

	return topAlbumsResp.TopAlbums.Albums
}

func FetchImage(imageUrl string, directory string, wg *sync.WaitGroup) (string, error) {
	defer wg.Done()

	resp, err := http.Get(imageUrl)
	if err != nil {
		slog.Error("Error fetching", "url", imageUrl, "error", err)
		return "", err
	}

	defer resp.Body.Close()

	filename := path.Join(directory, path.Base(resp.Request.URL.Path))
	out, err := os.Create(filename)
	if err != nil {
		slog.Error("Error creating image file", "error", err)
		return "", err
	}

	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		slog.Error("Error copying response image to created file", "error", err)
		return "", err
	}

	return filename, nil
}
