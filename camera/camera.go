package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"log"
	"os"
	"os/exec"
)

type Camera struct {
}

func NewCamera() *Camera {
	return &Camera{}
}

func (c *Camera) GetImage() (*image.Image, error) {
	return c.grabImageUsingCommand()
}

func (c *Camera) Close() {
}

func (c *Camera) grabImageUsingCommand() (*image.Image, error) {
	cmd := exec.Command("rpicam-still", "-n", "--immediate", "-o", "-")
	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %v -- %s", err, errBuffer.String())
	}
	imageBytes := outBuffer.Bytes()
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func main() {
	c := NewCamera()
	img, err := c.GetImage()
	if err != nil {
		panic(err)
	}
	f, err := os.Create("/tmp/output.jpg")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()
	options := &jpeg.Options{
		Quality: 85,
	}
	err = jpeg.Encode(f, *img, options)
	if err != nil {
		log.Fatalf("Failed to encode JPEG: %v", err)
	}
}
