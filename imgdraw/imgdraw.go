package imgdraw

import (
	"image"
	"image/color"
)

type ImageDrawer struct {
	img *image.RGBA
	col color.RGBA
}

func NewImageDrawer(img *image.RGBA) *ImageDrawer {
	return &ImageDrawer{img: img}
}

func (id *ImageDrawer) Color(col color.RGBA) {
	id.col = col
}

func (id *ImageDrawer) Point(x, y int) {
	id.img.SetRGBA(x, y, id.col)
}

func (id *ImageDrawer) Line(x0, y0, x1, y1 int) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	e := dx - dy
	for {
		id.Point(x0, y0)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * e
		if e2 > -dy {
			e -= dy
			x0 += sx
		}
		if e2 < dx {
			e += dx
			y0 += sy
		}
	}
}

func abs(n int) int {
	// Why is there no integer abs in Go?
	if n >= 0 {
		return n
	}
	return -n
}
