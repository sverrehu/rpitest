package motion_detector

import (
	"image"

	"golang.org/x/image/draw"
)

func MakeZeroBasedRGBA(img image.Image) *image.RGBA {
	if !needsNormalization(img) {
		return img.(*image.RGBA)
	}
	bounds := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(dst, dst.Bounds(), img, bounds.Min, draw.Src)
	return dst
}

func needsNormalization(img image.Image) bool {
	return !isRGBA(img) || img.Bounds().Min.X != 0 || img.Bounds().Min.Y != 0
}

func isRGBA(img image.Image) bool {
	switch img.(type) {
	case *image.RGBA:
		return true
	default:
		return false
	}
}
