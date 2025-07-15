package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/geavenx/smnh/cmd"
)

var (
	username  string
	rows      string
	columns   string
	artist    bool
	album     bool
	playcount bool
	method    string
	period    string

	prompt bool
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&username, "u", "", "Last.fm username (required)")
	flag.StringVar(&rows, "rows", "5", "Number of rows to display")
	flag.StringVar(&columns, "cols", "5", "Number of columns to display")
	flag.BoolVar(&artist, "artist", true, "Display artist")
	flag.BoolVar(&album, "album", true, "Display album")
	flag.BoolVar(&playcount, "playcount", false, "Display playcount")
	flag.StringVar(&method, "method", "album", "[album, artist, track]")
	flag.StringVar(&period, "period", "7day", "[overall, 7day, 1month, 6month, 12month]")
	flag.BoolVar(&prompt, "p", false, "Use prompt to define collage settings")

	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Fprintln(flag.CommandLine.Output(), "unexpected positional arguments")
		flag.Usage()
		os.Exit(2)
	}

	if prompt {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("username").
					Value(&username).
					Validate(func(s string) error {
						if s == "" {
							return errors.New("username cannot be empty")
						}
						return nil
					}),
				huh.NewInput().
					Title("rows").
					Value(&rows).
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
					Value(&columns).
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
					Value(&artist),
				huh.NewConfirm().
					Title("Display Album").
					Value(&album),
				huh.NewConfirm().
					Title("Display Playcount").
					Value(&playcount),

				huh.NewSelect[string]().Title("method").
					Options(
						huh.NewOption("Album", "album"),
						huh.NewOption("Artist", "artist"),
						huh.NewOption("Track", "track"),
					).Value(&method),
				huh.NewSelect[string]().Title("period").
					Options(
						huh.NewOption("Overall", "overall"),
						huh.NewOption("7 Day", "7day"),
						huh.NewOption("1 Month", "1month"),
						huh.NewOption("6 Month", "6month"),
						huh.NewOption("1 Year", "12month"),
					).Value(&period),
			),
		)

		err := form.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	cmd.Request(cmd.CollageRequest{Rows: rows, Columns: columns, Artist: artist, Playcount: playcount, Username: username, Period: period, Method: method, Album: album})
}
