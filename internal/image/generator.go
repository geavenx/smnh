package images

import (
	"crypto/rand"
	"log/slog"
	"path"

	"github.com/fogleman/gg"
)

func GenerateNotFoundImg(directory string) (string, error) {
	slog.Debug("Entering images.GenerateNotFoundImg() function")
	const W, H = 300, 300

	dc := gg.NewContext(W, H)

	var b [1]byte
	if _, err := rand.Read(b[:]); err != nil {
		slog.Error("Failed Generating random byte for not found image filename", "error", err)
		return "", err
	}
	filename := path.Join(directory, string(b[0]))
	slog.Debug("Generated filepath for not found image", "filepath", filename)

	dc.SetRGB(0, 0, 0)
	dc.Clear()

	if err := dc.LoadFontFace("assets/NotoSans-ExtraBold.ttf", 24); err != nil {
		slog.Error("Error loading font face", "filename", filename, "error", err)
		return "", err
	}

	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored("Cover not found", W/2, H/2, 0.5, 0.5)

	if err := dc.SavePNG(filename); err != nil {
		slog.Error("Failed saving not found image to filesystem", "filename", filename, "error", err)
		return "", err
	}

	return filename, nil
}
