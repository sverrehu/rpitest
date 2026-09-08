package main

// On my Mac: PKG_CONFIG_PATH=/opt/local/lib/opencv4/pkgconfig go run viewvideo.go

import (
	"gocv.io/x/gocv"
)

func main() {
	webcam, _ := gocv.OpenVideoCapture(0)
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		window.IMShow(img)
		window.WaitKey(1)
	}
}
