package motion_detector

import (
	"image"
	"image/color"
)

type Convolver struct {
	kernel       [][]float64
	kernelHeight int
	kernelWidth  int
}

func New3x3BlurConvolver() *Convolver {
	return NewConvolver([][]float64{
		{1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0},
		{1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0},
		{1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0},
	})
}

func NewConvolver(kernel [][]float64) *Convolver {
	kernelHeight := len(kernel)
	kernelWidth := len(kernel[0])
	return &Convolver{kernel: kernel, kernelHeight: kernelHeight, kernelWidth: kernelWidth}
}

func (c *Convolver) Convolve(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	for y := dst.Bounds().Max.Y - 1; y >= dst.Bounds().Min.Y; y-- {
		for x := dst.Bounds().Max.X - 1; x >= dst.Bounds().Min.X; x-- {
			dst.Set(x, y, c.convolve(src, x, y))
		}
	}
	return dst
}

func (c *Convolver) convolve(src image.Image, x int, y int) color.Color {
	rSum := 0.0
	gSum := 0.0
	bSum := 0.0
	for ky := c.kernelHeight - 1; ky >= 0; ky-- {
		kyy := y - c.kernelHeight/2 + ky
		if kyy < src.Bounds().Min.Y {
			kyy = src.Bounds().Min.Y
		}
		if kyy > src.Bounds().Max.Y {
			kyy = src.Bounds().Max.Y
		}
		for kx := c.kernelWidth - 1; kx >= 0; kx-- {
			kxx := x - c.kernelWidth/2 + kx
			if kxx < src.Bounds().Min.X {
				kxx = src.Bounds().Min.X
			}
			if kxx > src.Bounds().Max.X {
				kxx = src.Bounds().Max.X
			}
			weight := c.kernel[ky][kx]
			r, g, b, _ := src.At(kxx, kyy).RGBA()
			rSum += float64(r) * weight
			gSum += float64(g) * weight
			bSum += float64(b) * weight
		}
	}
	return color.RGBA64{R: uint16(rSum), G: uint16(gSum), B: uint16(bSum)}
}
