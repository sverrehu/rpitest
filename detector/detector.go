package detector

import (
	"image"
	"math"

	"golang.org/x/image/draw"
)

type Detector interface {
	Init() error
	Close()
	Detect(img *image.RGBA) ([]*Detection, error)
}

type Detection struct {
	X0, Y0 int
	X1, Y1 int
	Class  int
	Score  float32
}

func scaleImage(img *image.RGBA, width, height int) (*image.RGBA, float32) {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	hScale := float64(img.Bounds().Dx()) / float64(width)
	vScale := float64(img.Bounds().Dy()) / float64(height)
	scale := float32(math.Max(hScale, vScale))
	dstRect := image.Rect(0, 0, int(float32(img.Bounds().Dx())/scale), int(float32(img.Bounds().Dy())/scale))
	draw.ApproxBiLinear.Scale(dst, dstRect, img, img.Bounds(), draw.Over, nil)
	return dst, scale
}
