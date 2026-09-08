package main

// On my Mac: PKG_CONFIG_PATH=/opt/local/lib/opencv4/pkgconfig go run viewvideo.go

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"gocv.io/x/gocv"
)

func main() {
	webcam, err := gocv.OpenVideoCaptureWithAPI(0, gocv.VideoCaptureV4L2)
	if err != nil {
		log.Fatalf("Error opening webcam: %v", err)
	}
	defer webcam.Close()
	img := gocv.NewMat()
	defer img.Close()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	fmt.Println("Starting continuous image capture... Press Ctrl+C to stop.")
	for range ticker.C {
		if ok := webcam.Read(&img); !ok || img.Empty() {
			log.Println("Warning: Unable to read frame from webcam, skipping...")
			continue
		}
		buf, err := gocv.IMEncode(".jpg", img)
		if err != nil {
			log.Panic(err)
		}
		defer buf.Close()
		err = postImage(buf.GetBytes())
		if err != nil {
			log.Panic(err)
		}
	}
}

func postImage(b []byte) error {
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
