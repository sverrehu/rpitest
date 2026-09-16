package imgposter

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"time"
)

type ImagePoster struct {
	rpiviewURL string
}

func NewImagePoster(rpiviewHost string, rpiviewPort int) *ImagePoster {
	return &ImagePoster{rpiviewURL: fmt.Sprintf("http://%s:%d/img", rpiviewHost, rpiviewPort)}
}

func (ip *ImagePoster) PostJPEG(img image.Image) error {
	b, err := getJPEGBytes(img)
	if err != nil {
		return err
	}
	return ip.PostImageBytes(b, "image/jpeg")
}

func (ip *ImagePoster) PostImageBytes(b []byte, contentType string) error {
	r := bytes.NewReader(b)
	url := ip.rpiviewURL
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Post(url, contentType, r)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func getJPEGBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	options := &jpeg.Options{
		Quality: 85,
	}
	err := jpeg.Encode(&buf, img, options)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
