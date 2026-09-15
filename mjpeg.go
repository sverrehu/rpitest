package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/sverrehu/rpitest/camera"
	"github.com/sverrehu/rpitest/imgposter"
)

func main() {

	rpiview := imgposter.NewImagePoster("192.168.1.15", 8086)
	filename := "/Users/sverrehu/Dropbox/tmp/film.mjpeg"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var reader io.Reader = file
	n := 0
	camera.NewMJPEGSplitter(reader, func(m *camera.MJPEGSplitter) {
		n++
		log.Printf("New frame %d", n)
		rpiview.PostImageBytes(m.GetLastImageBytes(), "image/jpeg")
	})
	time.Sleep(10 * time.Second) // to let the go routine run
}
