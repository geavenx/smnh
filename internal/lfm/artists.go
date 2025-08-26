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

type TopArtistRequest struct {
	Username string
	Period   Period
	Limit    int
}

type TopArtistResponse struct {
	XMLName    xml.Name  `xml:"lfm"`
	Status     string    `xml:"status,attr"`
	TopArtists TopArtist `xml:"topartists"`
}

type TopArtist struct {
	XMLName xml.Name `xml:"topartists"`
	User    string   `xml:"user,attr"`
	Artists []Artist `xml:"artist"`
}

type Artist struct {
	XMLName   xml.Name `xml:"artist"`
	Rank      string   `xml:"rank,attr"`
	Name      string   `xml:"name"`
	PlayCount int      `xml:"playcount"`
	MBID      string   `xml:"mbid"`
	Url       string   `xml:"url"`
	Images    []Image  `xml:"image"`
}

func FetchTopArtists(req TopArtistRequest) []Artist {
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
		"method":  "user.getTopArtists",
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

	slog.Debug("Sending request", "method", "user.getTopArtists", "url", apiUrl.String())
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

	var topArtistsResp TopArtistResponse
	if err = xml.Unmarshal(body, &topArtistsResp); err != nil {
		slog.Error("Failed to parse XML", "error", err)
	}

	return topArtistsResp.TopArtists.Artists
}
