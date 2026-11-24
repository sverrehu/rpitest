package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/sverrehu/rpitest/camera"
	"github.com/sverrehu/rpitest/motion_detector"
)

const outputDir = "/home/pi/tmp"
const timestampFormat = "2006-01-02.150405.000"

func main() {
	//rpiviewHost := "192.168.1.15"
	////rpiviewHost := "192.168.30.21"
	//rpiview := imgposter.NewImagePoster(rpiviewHost, 8086)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(err)
	}
	md := motion_detector.NewDefaultMotionDetector()
	cam := camera.NewCamera()
	defer cam.Close()
	cam.Rotate = true
	err = cam.StartStreaming(nil)
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(250 * time.Millisecond)
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
		if md.IsMotion(img) {
			//fmt.Println("motion")
			//err = rpiview.PostJPEG(md.GetPrevImg()) //img)
			err = saveMotionImage(img)
			if err != nil {
				log.Panic(err)
			}
		}
	}
}

func saveMotionImage(img image.Image) error {
	filename := outputDir + "/motion-" + time.Now().Format(timestampFormat) + ".jpg"
	return saveJPEG(filename, img)
}

func saveJPEG(filename string, img image.Image) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, img, &jpeg.Options{Quality: 85})
}
