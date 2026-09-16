package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sverrehu/rpitest/camera"
	"github.com/sverrehu/rpitest/imgposter"
)

func main() {
	rpiview := imgposter.NewImagePoster("192.168.1.15", 8086)
	webcam := camera.NewCamera()
	defer webcam.Close()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	fmt.Println("Starting continuous image capture... Press Ctrl+C to stop.")
	for range ticker.C {
		startTime := time.Now()
		img, err := webcam.GetSingleImage()
		elapsed := time.Since(startTime)
		log.Printf("Captured image in %v", elapsed)
		if err != nil {
			log.Panic(err)
		}
		err = rpiview.PostJPEG(img)
		if err != nil {
			log.Panic(err)
		}
	}
}
