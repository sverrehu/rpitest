package main

// On my Mac: PKG_CONFIG_PATH=/opt/local/lib/opencv4/pkgconfig go run viewvideo.go

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"time"

	"github.com/sverrehu/rpitest/camera"
)

func main() {
	webcam := camera.NewCamera()
	defer webcam.Close()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	fmt.Println("Starting continuous image capture... Press Ctrl+C to stop.")
	for range ticker.C {
		startTime := time.Now()
		img, err := webcam.GetImage()
		elapsed := time.Since(startTime)
		log.Printf("Captured image in %v", elapsed)
		if err != nil {
			log.Panic(err)
		}
		imgBytes, err := getImageBytes(img)
		if err != nil {
			log.Panic(err)
		}
		err = postImageBytes(imgBytes)
		if err != nil {
			log.Panic(err)
		}
	}
}

func getImageBytes(img *image.Image) ([]byte, error) {
	var buf bytes.Buffer
	options := &jpeg.Options{
		Quality: 85,
	}
	err := jpeg.Encode(&buf, *img, options)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func postImageBytes(b []byte) error {
	r := bytes.NewReader(b)
	url := "http://192.168.1.15:8086/img"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Post(url, "image/jpeg", r)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}
