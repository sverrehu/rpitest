package motion_detector

import (
	"fmt"
	"image"

	"github.com/sverrehu/goutils/image/scaling"
	"golang.org/x/image/draw"
)

type MotionDetector struct {
	noiseThreshold           int
	requiredChangePercentage float64
	scaling                  int
	prevImg                  *image.Gray // converted to gray and scaled
}

func NewDefaultMotionDetector() *MotionDetector {
	return NewMotionDetector(12, 1.5, 8)
}

func NewMotionDetector(noiseThreshold int, requiredChangePercentage float64, scaling int) *MotionDetector {
	return &MotionDetector{noiseThreshold, requiredChangePercentage, scaling, nil}
}

func (m *MotionDetector) IsMotion(img image.Image) bool {
	if img == nil {
		return false
	}
	newImg := m.prepare(img)
	if m.prevImg == nil {
		m.prevImg = newImg
		return false
	}
	if m.prevImg.Bounds().Dx() != newImg.Bounds().Dx() || m.prevImg.Bounds().Dy() != newImg.Bounds().Dy() {
		m.prevImg = newImg
		return true
	}
	numSamples := 0
	numChanged := 0
	bounds := newImg.Bounds()
	// newImg is zero-based grayscale
	for y := bounds.Max.Y - 1; y >= 0; y-- {
		for x := bounds.Max.X - 1; x >= 0; x-- {
			oldGray := m.prevImg.Pix[(y*m.prevImg.Stride)+x]
			newGray := newImg.Pix[(y*newImg.Stride)+x]
			absDiff := abs(int(oldGray) - int(newGray))
			if absDiff > m.noiseThreshold {
				//fmt.Println("absDiff", absDiff)
				numChanged++
			}
			numSamples++
		}
	}
	changedPercentage := 100.0 * float64(numChanged) / float64(numSamples)
	fmt.Println("changedPercentage", changedPercentage)
	isMotion := changedPercentage > m.requiredChangePercentage
	if isMotion {
		m.prevImg = newImg
	}
	return isMotion
}

func (m *MotionDetector) GetPrevImg() *image.Gray {
	return m.prevImg
}

func (m *MotionDetector) prepare(img image.Image) *image.Gray {
	dst := img
	if m.scaling > 1 {
		dst = scaling.Resize(img, img.Bounds().Dx()/m.scaling, img.Bounds().Dy()/m.scaling)
	}
	bounds := dst.Bounds()
	gray := image.NewGray(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(gray, gray.Bounds(), dst, bounds.Min, draw.Src)
	return gray
}

func (m *MotionDetector) grayAt(img image.Image, x, y int) int {
	r, g, b, _ := img.At(x, y).RGBA()
	r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
	return int(r8+g8+b8) / 3
}

func abs(n int) int {
	// Why is there no integer abs in Go?
	if n >= 0 {
		return n
	}
	return -n
}
