package lfm

import (
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type LfmRequest struct {
	Method   Method `xml:"method"`
	Username string `xml:"username"`
	Period   Period `xml:"period"`
	Limit    int    `xml:"limit"`
}

type LfmObject struct {
	Type      string
	Name      string
	PlayCount int
	Rank      string
	Images    []Image
}

func FetchLfm(req LfmRequest) {
	if err := godotenv.Load(); err != nil {
		slog.Error("error loading .env", "error", err)
	}

	baseUrl := os.Getenv("LAST_FM_BASE_URL")
	apiKey := os.Getenv("LAST_FM_API_KEY")

	apiUrl, err := url.Parse(baseUrl)
	if err != nil {
		slog.Error("error parsing last.fm API URL", "error", err)
	}

	params := map[string]string{
		"method":  string(req.Method),
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

	slog.Debug("Sending request", "method", req.Method, "url", apiUrl.String())
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

	switch req.Method {
	case TopAlbums:
		var resp TopAlbumResponse
	case TopArtists:
		var resp TopArtistResponse
	default:
		var resp TopAlbumResponse
	}

	if err = xml.Unmarshal(body, &resp); err != nil {
		slog.Error("Failed to parse XML", "error", err)
	}

	return topAlbumsResp.TopAlbums.Albums
}

// func FetchImage(imageUrl string, directory string, wg *sync.WaitGroup) (string, error) {
// 	defer wg.Done()
//
// 	resp, err := http.Get(imageUrl)
// 	if err != nil {
// 		slog.Error("Error fetching", "url", imageUrl, "error", err)
// 		return "", err
// 	}
//
// 	defer resp.Body.Close()
//
// 	filename := path.Join(directory, path.Base(resp.Request.URL.Path))
// 	out, err := os.Create(filename)
// 	if err != nil {
// 		slog.Error("Error creating image file", "error", err)
// 		return "", err
// 	}
//
// 	defer out.Close()
//
// 	if _, err = io.Copy(out, resp.Body); err != nil {
// 		slog.Error("Error copying response image to created file", "error", err)
// 		return "", err
// 	}
//
// 	return filename, nil
// }
