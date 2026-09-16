package camera

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"log"
	"os/exec"
)

type Camera struct {
	listener      func(*Camera)
	mjpegSplitter *MJPEGSplitter
}

func NewCamera() *Camera {
	return &Camera{}
}

func (c *Camera) GetSingleImage() (*image.Image, error) {
	return c.grabSingleImageUsingCommand()
}

func (c *Camera) GetImage() (*image.Image, error) {
	if c.mjpegSplitter == nil {
		return c.GetSingleImage()
	}
	return c.mjpegSplitter.GetLastImage(), nil
}

func (c *Camera) StartStreaming(listener func(*Camera)) error {
	c.listener = listener
	return c.grabStreamUsingCommand()
}

func (c *Camera) Close() {
	log.Print("Closing camera")
	if c.mjpegSplitter != nil {
		c.mjpegSplitter.terminate = true
	}
}

func (c *Camera) grabSingleImageUsingCommand() (*image.Image, error) {
	cmd := exec.Command("rpicam-still", "--nopreview", "--zsl", "--immediate", "--thumb", "none", "--exposure", "sport", "-o", "-")
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

func (c *Camera) grabStreamUsingCommand() error {
	cmd := exec.Command("rpicam-vid", "--nopreview", "-t", "0", "--codec", "mjpeg", "--quality", "85", "--inline", "-o", "-")
	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer
	go func() {
		log.Print("Running rpicam-vid")
		err := cmd.Run()
		if err != nil {
			log.Fatalf("failed to execute command: %v -- %s", err, errBuffer.String())
		}
		log.Print("Done running rpicam-vid")
	}()
	c.mjpegSplitter = NewMJPEGSplitter(&outBuffer, func(splitter *MJPEGSplitter) {
		if c.listener != nil {
			c.listener(c)
		}
	})
	return nil
}
