package main

import (
	"fmt"
	"image"
	"log"

	"github.com/sverrehu/rpitest/camera"
	"github.com/sverrehu/rpitest/imgposter"
	"github.com/sverrehu/rpitest/motion_detector"
)

func main() {
	rpiviewHost := "192.168.1.15"
	//rpiviewHost := "192.168.30.21"
	rpiview := imgposter.NewImagePoster(rpiviewHost, 8086)
	md := motion_detector.NewDefaultMotionDetector()
	conv := motion_detector.New3x3BlurConvolver()
	cam := camera.NewCamera()
	defer cam.Close()
	cam.Rotate = true
	err := cam.StartStreaming(nil)
	if err != nil {
		panic(err)
	}
	oldImg := image.Image(nil)
	fmt.Println("Starting continuous image capture... Press Ctrl+C to stop.")
	for {
		img, err := cam.GetImage()
		if img == nil {
			continue
		}
		if err != nil {
			log.Panic(err)
		}
		blurred := conv.Convolve(img)
		if md.IsMotion(oldImg, blurred) {
			fmt.Println("motion")
			err = rpiview.PostJPEG(blurred)
			if err != nil {
				log.Panic(err)
			}
		}
		oldImg = blurred
	}
}
