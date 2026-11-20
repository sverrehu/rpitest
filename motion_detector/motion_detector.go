package motion_detector

import (
	"fmt"
	"image"
)

type MotionDetector struct {
	noiseThreshold           int
	requiredChangePercentage float64
}

func NewDefaultMotionDetector() *MotionDetector {
	return NewMotionDetector(80, 1.5)
}

func NewMotionDetector(noiseThreshold int, requiredChangePercentage float64) *MotionDetector {
	return &MotionDetector{noiseThreshold, requiredChangePercentage}
}

func (m *MotionDetector) IsMotion(oldImg image.Image, newImg image.Image) bool {
	if oldImg == nil || newImg == nil {
		return false
	}
	if oldImg.Bounds().Dx() != newImg.Bounds().Dx() || oldImg.Bounds().Dy() != newImg.Bounds().Dy() {
		return false
	}
	numSamples := 0
	numChanged := 0
	bounds := newImg.Bounds()
	for y := bounds.Max.Y - 1; y >= bounds.Min.Y; y -= 8 {
		for x := bounds.Max.X - 1; x >= bounds.Min.X; x -= 8 {
			oldGray := m.grayAt(oldImg, x, y)
			newGray := m.grayAt(newImg, x, y)
			absDiff := abs(oldGray - newGray)
			if absDiff > m.noiseThreshold {
				fmt.Println("absDiff", absDiff)
				numChanged++
			}
			numSamples++
		}
	}
	changedPercentage := 100.0 * float64(numChanged) / float64(numSamples)
	fmt.Println("changedPercentage", changedPercentage)
	return changedPercentage > m.requiredChangePercentage
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
