package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sverrehu/rpitest/camera"
	"github.com/sverrehu/rpitest/imgposter"
)

func main() {
	//rpiviewHost := "192.168.1.15"
	rpiviewHost := "192.168.30.21"
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
		err = rpiview.PostJPEG(img)
		if err != nil {
			log.Panic(err)
		}
	}
}
