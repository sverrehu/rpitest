package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/sverrehu/rpitest/camera"
	"github.com/sverrehu/rpitest/imgdraw"
	"github.com/sverrehu/rpitest/imgposter"
)

func main() {
	rpiviewHost := "192.168.1.15"
	//rpiviewHost := "192.168.30.21"
	rpiview := imgposter.NewImagePoster(rpiviewHost, 8086)
	cam := camera.NewCamera()
	defer cam.Close()
	cam.Rotate = true
	err := cam.StartStreaming(nil)
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	fmt.Println("Starting continuous image capture... Press Ctrl+C to stop.")
	for range ticker.C {
		img, err := cam.GetImage()
		if img == nil {
			continue
		}
		if err != nil {
			log.Panic(err)
		}
		annotate(img)
		err = rpiview.PostJPEG(img)
		if err != nil {
			log.Panic(err)
		}
	}
}

func annotate(img image.Image) {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		panic("Image is not RGBA")
	}
	id := imgdraw.NewImageDrawer(rgba)
	bounds := rgba.Bounds()
	id.Color(color.RGBA{255, 100, 100, 255})
	id.Line(0, bounds.Dy()/2, bounds.Dx(), bounds.Dy()/2)
	id.Line(bounds.Dx()/2, 0, bounds.Dx()/2, bounds.Dy())
}
