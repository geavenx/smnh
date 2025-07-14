package main

import (
	"errors"
	"fmt"
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
)

func main() {
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

	cmd.Request(cmd.CollageRequest{Rows: rows, Columns: columns, Artist: artist, Playcount: playcount, Username: username, Period: period, Method: method, Album: album})
}
