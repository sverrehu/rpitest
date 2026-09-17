package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/sverrehu/rpitest/camera"
	"github.com/sverrehu/rpitest/detector"
	"github.com/sverrehu/rpitest/imgdraw"
	"github.com/sverrehu/rpitest/imgposter"
)

func main() {
	rpiviewHost := "192.168.1.15"
	//rpiviewHost := "192.168.30.21"
	rpiview := imgposter.NewImagePoster(rpiviewHost, 8086)
	det := detector.NewONNXRuntimeDetector("../gotest/gocv/yolo26_face_fp16.onnx", 640, 640)
	err := det.Init()
	if err != nil {
		panic(err)
	}
	cam := camera.NewCamera()
	defer cam.Close()
	cam.Rotate = true
	err = cam.StartStreaming(nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting continuous image capture... Press Ctrl+C to stop.")
	for {
		img, err := cam.GetImage()
		if img == nil {
			continue
		}
		if err != nil {
			log.Panic(err)
		}
		detections, err := det.Detect(img)
		if err != nil {
			panic(err)
		}
		annotateDetections(img, detections)
		err = rpiview.PostJPEG(img)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Image size: %d x %d", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func annotateDetections(img image.Image, detections []*detector.Detection) {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		panic("Image is not RGBA")
	}
	id := imgdraw.NewImageDrawer(rgba)
	bounds := rgba.Bounds()
	id.Color(color.RGBA{255, 100, 100, 255})
	id.Line(0, bounds.Dy()/2, bounds.Dx(), bounds.Dy()/2)
	id.Line(bounds.Dx()/2, 0, bounds.Dx()/2, bounds.Dy())
	id.Color(color.RGBA{100, 255, 100, 255})
	for _, d := range detections {
		id.Rect(d.X0, d.Y0, d.X1, d.Y1)
	}
}
