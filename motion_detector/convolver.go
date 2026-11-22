package motion_detector

import (
	"image"
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

func New5x5BlurConvolver() *Convolver {
	return NewConvolver([][]float64{
		{1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0},
		{1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0},
		{1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0},
		{1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0},
		{1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0, 1.0 / 25.0},
	})
}

func NewConvolver(kernel [][]float64) *Convolver {
	kernelHeight := len(kernel)
	kernelWidth := len(kernel[0])
	return &Convolver{kernel: kernel, kernelHeight: kernelHeight, kernelWidth: kernelWidth}
}

func (c *Convolver) Convolve(img image.Image) *image.RGBA {
	src := MakeZeroBasedRGBA(img)
	dst := image.NewRGBA(src.Bounds())
	bounds := dst.Bounds()
	srcPix := src.Pix
	dstPix := dst.Pix
	srcStride := src.Stride
	dstStride := dst.Stride
	// both src and dest are zero-based RGBA images with the same bounds
	for y := bounds.Max.Y - 1; y >= 0; y-- {
		for x := bounds.Max.X - 1; x >= 0; x-- {
			dstIdx := (y * dstStride) + (x * 4)
			r, g, b := c.convolve(srcPix, srcStride, bounds.Max.X-1, bounds.Max.Y-1, x, y)
			dstPix[dstIdx] = r
			dstPix[dstIdx+1] = g
			dstPix[dstIdx+2] = b
		}
	}
	return dst
}

func (c *Convolver) convolve(srcPix []uint8, srcStride int, srcMaxX int, srcMaxY int, x int, y int) (uint8, uint8, uint8) {
	// expects img to be normalized
	rSum := 0.0
	gSum := 0.0
	bSum := 0.0
	for ky := c.kernelHeight - 1; ky >= 0; ky-- {
		kyy := y - c.kernelHeight/2 + ky
		if kyy < 0 {
			kyy = 0
		}
		if kyy > srcMaxY {
			kyy = srcMaxY
		}
		for kx := c.kernelWidth - 1; kx >= 0; kx-- {
			kxx := x - c.kernelWidth/2 + kx
			if kxx < 0 {
				kxx = 0
			}
			if kxx > srcMaxX {
				kxx = srcMaxX
			}
			weight := c.kernel[ky][kx]
			srcIdx := (kyy * srcStride) + (kxx * 4)
			//fmt.Printf("kxx: %d, kyy: %d, weight: %f, srcIdx: %d\n", kxx, kyy, weight, srcIdx)
			r := srcPix[srcIdx]
			g := srcPix[srcIdx+1]
			b := srcPix[srcIdx+2]
			rSum += float64(r) * weight
			gSum += float64(g) * weight
			bSum += float64(b) * weight
		}
	}
	return uint8(rSum), uint8(gSum), uint8(bSum)
}
