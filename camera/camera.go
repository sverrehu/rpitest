package camera

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"log"
	"os/exec"
)

type Camera struct {
	Rotate        bool
	listener      func(*Camera)
	mjpegSplitter *MJPEGSplitter
}

func NewCamera() *Camera {
	return &Camera{}
}

func (c *Camera) GetSingleImage() (*image.RGBA, error) {
	return c.grabSingleImageUsingCommand()
}

func (c *Camera) GetImage() (*image.RGBA, error) {
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

func (c *Camera) grabSingleImageUsingCommand() (*image.RGBA, error) {
	args := []string{"--nopreview", "--zsl", "--immediate", "--thumb", "none", "--exposure", "sport", "-o", "-"}
	if c.Rotate {
		args = append(args, "--rotation", "180")
	}
	cmd := exec.Command("rpicam-still", args...)
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
	rgba := toRGBA(img)
	return rgba, nil
}

func (c *Camera) grabStreamUsingCommand() error {
	args := []string{"--nopreview", "-t", "0", "--codec", "mjpeg", "--quality", "85", "--inline", "-o", "-"}
	if c.Rotate {
		args = append(args, "--rotation", "180")
	}
	cmd := exec.Command("rpicam-vid", args...)
	var errBuffer bytes.Buffer
	cmd.Stderr = &errBuffer
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go func() {
		log.Print("Running rpicam-vid")
		err := cmd.Run()
		if err != nil {
			log.Fatalf("failed to execute command: %v -- %s", err, errBuffer.String())
		}
		log.Print("Done running rpicam-vid")
	}()
	c.mjpegSplitter = NewMJPEGSplitter(outPipe, func(splitter *MJPEGSplitter) {
		if c.listener != nil {
			c.listener(c)
		}
	})
	return nil
}

func toRGBA(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}
	b := img.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)
	return rgba
}
