package main

import (
	"io"
	"os"

	"github.com/sverrehu/rpitest/camera"
)

func main() {
	filename := "/Users/sverrehu/Dropbox/tmp/film.mjpeg"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var reader io.Reader = file
	_ = camera.NewMJPEGSplitter(reader, func() {
		println("New frame")
	})
}
