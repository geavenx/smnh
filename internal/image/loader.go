package images

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log/slog"
	"os"
	"slices"
	"sync"
)

func LoadImage(filename string, wg *sync.WaitGroup) (image.Image, error) {
	defer wg.Done()

	file, err := os.Open(filename)
	if err != nil {
		slog.Error("Failed to open image file", "filename", filename, "error", err)
		return nil, err
	}
	defer file.Close()

	// Automatic format detection
	img, format, err := image.Decode(file)
	if err != nil {
		slog.Error("Failed to decode image", "filename", filename, "error", err)
		return nil, err
	}

	// Validade supported formats
	supportedFormats := []string{"jpeg", "jpg", "png", "gif"}
	if !slices.Contains(supportedFormats, format) {
		slog.Error("Unsupported format", "format", format)
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	return img, nil
}

// func ExtractFirstFrame(gifPath string) (image.Image, error) {
// 	file, err := os.Open(gifPath)
// 	if err != nil {
// 		slog.Error("Failed to open GIF file", "filepath", gifPath, "error", err)
// 		return nil, err
// 	}
// 	defer file.Close()
//
// 	img, err := gif.Decode(file) // Decode only the first frame with gif.Decode()
// 	if err != nil {
// 		slog.Error("Failed to decode gif file", "filepath", gifPath, "error", err)
// 		return nil, err
// 	}
//
// 	return img, nil
// }

func SaveImage(img image.Image, filename string, format string) error {
	file, err := os.Create(filename)
	if err != nil {
		slog.Error("Failed to create file", "filename", filename, "error", err)
		return err
	}
	defer file.Close()

	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	case "png":
		return png.Encode(file, img)
	case "gif":
		return gif.Encode(file, img, nil)
	default:
		slog.Error("Unsupported output file format", "format", format)
		os.Exit(2)
		return fmt.Errorf("Unsupported format: %s", format)
	}
}

// ResizeImage resizes an image to the specified width and height
func ResizeImage(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	if bounds.Dx() == width && bounds.Dy() == height {
		return img // Already the correct size
	}

	// Create a new RGBA image with the target dimensions
	resized := image.NewRGBA(image.Rect(0, 0, width, height))

	// Simple nearest neighbor scaling
	scaleX := float64(bounds.Dx()) / float64(width)
	scaleY := float64(bounds.Dy()) / float64(height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := int(float64(x) * scaleX)
			srcY := int(float64(y) * scaleY)

			// Ensure we don't go out of bounds
			if srcX >= bounds.Dx() {
				srcX = bounds.Dx() - 1
			}
			if srcY >= bounds.Dy() {
				srcY = bounds.Dy() - 1
			}

			color := img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY)
			resized.Set(x, y, color)
		}
	}

	return resized
}
