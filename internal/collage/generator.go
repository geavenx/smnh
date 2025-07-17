package collage

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log/slog"
	"sync"
)

type CollageGenerator struct {
	config Config
	images []image.Image
	mutex  sync.Mutex
}

type Config struct {
	Width   int
	Height  int
	Rows    int
	Columns int
	Quality int
}

func NewCollageGenerator(config Config) *CollageGenerator {
	return &CollageGenerator{
		config: config,
		images: make([]image.Image, 0),
		mutex:  sync.Mutex{},
	}
}

func (g *CollageGenerator) AddImage(img image.Image) error {
	if img == nil {
		slog.Error("Failed because image is null")
		return fmt.Errorf("Image cannot be null")
	}

	// validate dimensions
	bounds := img.Bounds()
	if bounds.Dx() != 300 || bounds.Dy() != 300 {
		slog.Error("Images must be 300x300", "x_dimension", bounds.Dx(), "y_dimension", bounds.Dy())
		return fmt.Errorf("image must be 300x300 pixel, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.images = append(g.images, img)
	return nil
}

func (g *CollageGenerator) Generate() (image.Image, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if len(g.images) == 0 {
		slog.Error("CollageGenerator.Generate() called with no images to process")
		return nil, fmt.Errorf("No images to process")
	}
	return g.createCollage()
}

func (g *CollageGenerator) createCollage() (*image.RGBA, error) {
	canvas := createImageGrid(g.config.Rows, g.config.Columns, 300, g.images)

	return canvas, nil
}

// func calculateCanvasSize(cols, rows, cellWidth, cellHeight int) (int, int) {
// 	width := cols * cellWidth
// 	height := rows * cellHeight
//
// 	return width, height
// }

func createImageGrid(rows, cols, cellSize int, images []image.Image) *image.RGBA {
	canvasWidth := cols * cellSize
	canvasHeight := rows * cellSize

	canvas := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))

	bgColor := color.RGBA{255, 255, 255, 255}
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	for i, img := range images {
		if i >= rows*cols {
			break
		}

		row := i / cols
		col := i % cols

		x := col * cellSize
		y := row * cellSize

		destRect := image.Rect(x, y, x+cellSize, y+cellSize)
		draw.Draw(canvas, destRect, img, img.Bounds().Min, draw.Over)

	}

	return canvas

}
