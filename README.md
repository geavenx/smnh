# smnh

A CLI tool to generate music collages from your Last.fm scrobbling history.

## Features

- Generate album art collages from your Last.fm top albums
- Interactive prompt mode or command-line flags
- Customizable grid size (rows/columns)
- Multiple time periods (7 days, 1 month, 6 months, 1 year, overall)

## Installation

```bash
go build -o smnh .
```

## Usage

### Interactive Mode

```bash
./smnh -p
```

### Command Line Mode

```bash
./smnh -u your_lastfm_username -rows 5 -cols 5 -period 7day
```

### Options

- `-u`: Last.fm username (required)
- `-rows`: Number of rows (default: 5)
- `-cols`: Number of columns (default: 5)
- (**WIP**) `-method`: Data type [album, artist, track] (default: album)
- `-period`: Time period [overall, 7day, 1month, 6month, 12month] (default: 7day)
- (**WIP**) `-artist`: Display artist name (default: true)
- (**WIP**) `-album`: Display album name (default: true)
- (**WIP**) `-playcount`: Display play count (default: false)
- `-p`: Use interactive prompt
- `-debug`: Enable debug logging

## Output

The tool generates a PNG collage saved as `{username}_collage.png` in the current
directory.

## Requirements

- Go 1.24.5+
- Last.fm API Key and Shared Secret

