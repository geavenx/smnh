package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/charmbracelet/huh"
	"github.com/geavenx/smnh/internal/collage"
	images "github.com/geavenx/smnh/internal/image"
	"github.com/geavenx/smnh/internal/lfm"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags]\n", os.Args[0])
		flag.PrintDefaults()
	}

	var (
		username  = flag.String("u", "", "Last.fm username (required)")
		rows      = flag.String("rows", "5", "Number of rows to display")
		columns   = flag.String("cols", "5", "Number of columns to display")
		artist    = flag.Bool("artist", true, "Display artist")
		album     = flag.Bool("album", true, "Display album")
		playcount = flag.Bool("playcount", false, "Display playcount")
		method    = flag.String("method", "album", "[album, artist, track]")
		period    = flag.String("period", "7day", "[overall, 7day, 1month, 6month, 12month]")
		prompt    = flag.Bool("p", false, "Use prompt to define collage settings")
		debug     = flag.Bool("debug", false, "Enable debug logging")
	)
	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Fprintln(flag.CommandLine.Output(), "unexpected positional arguments")
		flag.Usage()
		os.Exit(2)
	}

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if *debug {
		opts.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	if *prompt {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("username").
					Value(username).
					Validate(func(s string) error {
						if s == "" {
							return errors.New("username cannot be empty")
						}
						return nil
					}),
				huh.NewInput().
					Title("rows").
					Value(rows).
					Prompt("[1 - 20]: ").
					Validate(func(s string) error {
						if s == "" {
							return errors.New("a row number is required")
						}
						num, err := strconv.Atoi(s)
						if err != nil {
							return errors.New("invalid row number")
						}

						if num < 1 || num > 20 {
							return errors.New("row number must be between 1 and 20")
						}

						return nil
					}),
				huh.NewInput().
					Title("columns").
					Value(columns).
					Prompt("[1 - 20]: ").
					Validate(func(s string) error {
						if s == "" {
							return errors.New("a column number is required")
						}
						num, err := strconv.Atoi(s)
						if err != nil {
							return errors.New("invalid column number")
						}

						if num < 1 || num > 20 {
							return errors.New("column number must be between 1 and 20")
						}

						return nil
					}),

				huh.NewConfirm().
					Title("Display Artist").
					Value(artist),
				huh.NewConfirm().
					Title("Display Album").
					Value(album),
				huh.NewConfirm().
					Title("Display Playcount").
					Value(playcount),

				huh.NewSelect[string]().Title("method").
					Options(
						huh.NewOption("Album", "album"),
						huh.NewOption("Artist", "artist"),
						huh.NewOption("Track", "track"),
					).Value(method),
				huh.NewSelect[string]().Title("period").
					Options(
						huh.NewOption("Overall", "overall"),
						huh.NewOption("7 Day", "7day"),
						huh.NewOption("1 Month", "1month"),
						huh.NewOption("6 Month", "6month"),
						huh.NewOption("1 Year", "12month"),
					).Value(period),
			),
		)

		err := form.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Convert period string to Period type
	var periodType lfm.Period
	switch *period {
	case "overall":
		periodType = lfm.PeriodOverall
	case "7day":
		periodType = lfm.Period7Day
	case "1month":
		periodType = lfm.Period1Month
	case "3month":
		periodType = lfm.Period3Month
	case "6month":
		periodType = lfm.Period6Month
	case "12month":
		periodType = lfm.Period12Month
	default:
		periodType = lfm.Period7Day
	}

	rowsInt, _ := strconv.Atoi(*rows)
	colsInt, _ := strconv.Atoi(*columns)
	totalImages := rowsInt * colsInt

	switch *method {
	case "album":
		res := lfm.RequestLfm(lfm.LfmRequest{Method: "user.getTopAlbums", Username: *username, Period: periodType, Limit: totalImages})
		top = lfm.FetchTopAlbums(lfm.TopAlbumRequest{Username: *username, Period: periodType, Limit: totalImages})
		slog.Debug("fetch top DONE", "top", top)
	case "artist":
		top = lfm.FetchTopArtists(lfm.TopArtistRequest{Username: *username, Period: periodType, Limit: totalImages})
		slog.Debug("fetch topArtists DONE", "top", top)
	default:
		top = lfm.FetchTopAlbums(lfm.TopAlbumRequest{Username: *username, Period: periodType, Limit: totalImages})
		slog.Debug("fetch top DONE", "top", top)
	}

	// Create temp directory for images
	dir, err := os.MkdirTemp("", "images")
	if err != nil {
		slog.Error("Error creating temp directory", "error", err)
	}
	defer os.RemoveAll(dir)

	var wg sync.WaitGroup
	results := make(chan lfm.ImageResult, len(top))

	for _, album := range top {
		wg.Add(1)

		// For each album in the top range run an asynchronous web request to download the image file to temp dir

		go func(url string, dir string) {
			filename, err := lfm.FetchImage(url, dir, &wg)
			results <- lfm.ImageResult{Filename: filename, Err: err}
		}(album.Images[3].Url, dir)
	}
	wg.Wait() // Wait for all goroutines to be done

	close(results) // Close the channel
	slog.Info("done")

	generator := collage.NewCollageGenerator(collage.Config{
		Width:   colsInt * 300,
		Height:  rowsInt * 300,
		Rows:    rowsInt,
		Columns: colsInt,
		Quality: 90,
	})

	// Collect successful image filenames
	var imageFilenames []string
	for res := range results {
		if res.Err != nil {
			slog.Warn("Error fetching image", "error", res.Err, "filename", res.Filename)
			filename, err := images.GenerateNotFoundImg(dir)
			if err != nil {
				slog.Error("Error generating 'not found' image", "error", err)
			} else {
				slog.Info("Not found image generated succesfully", "filename", filename)
				imageFilenames = append(imageFilenames, filename)
			}
		} else {
			slog.Info("Image fetched", "filename", res.Filename)
			imageFilenames = append(imageFilenames, res.Filename)
		}
	}

	// Load images and add to collage generator
	var imageWg sync.WaitGroup
	for _, filename := range imageFilenames {
		imageWg.Add(1)
		go func(filename string) {
			img, err := images.LoadImage(filename, &imageWg)
			if err != nil {
				slog.Error("Error loading image", "filename", filename, "error", err)
				return
			}

			// Resize image to 300x300 for collage
			resizedImg := images.ResizeImage(img, 300, 300)

			err = generator.AddImage(resizedImg)
			if err != nil {
				slog.Error("Error adding image to collage", "filename", filename, "error", err)
			}
		}(filename)
	}

	imageWg.Wait()

	// Generate the collage
	slog.Info("Generating collage...")
	collageImg, err := generator.Generate()
	if err != nil {
		slog.Error("Error generating collage", "error", err)
		os.Exit(2)
	}

	// Save the collage
	outputFilename := fmt.Sprintf("%s_collage.png", *username)
	err = images.SaveImage(collageImg, outputFilename, "png")
	if err != nil {
		slog.Error("Error saving collage", "filename", outputFilename, "error", err)
		os.Exit(2)
	}

	slog.Info("Collage generated successfully", "filename", outputFilename)
	fmt.Printf("Collage saved as: %s\n", outputFilename)
}
