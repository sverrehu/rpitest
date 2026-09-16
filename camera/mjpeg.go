package camera

import (
	"bytes"
	"image"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"io"
	"log"
	"sync"
)

var (
	soiMarker = []byte{0xff, 0xd8}
	eoiMarker = []byte{0xff, 0xd9}
)

type MJPEGSplitter struct {
	stream         io.Reader
	terminate      bool
	lastImageBytes []byte
	listener       func(*MJPEGSplitter)
	lastImageMutex sync.Mutex
}

func NewMJPEGSplitter(stream io.Reader, listener func(*MJPEGSplitter)) *MJPEGSplitter {
	m := &MJPEGSplitter{stream: stream, terminate: false, listener: listener}
	go m.inputHandlerLoop()
	return m
}

func (m *MJPEGSplitter) GetLastImage() *image.Image {
	imageBytes := m.GetLastImageBytes()
	if imageBytes == nil {
		return nil
	}
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil
	}
	return &img
}

func (m *MJPEGSplitter) GetLastImageBytes() []byte {
	m.lastImageMutex.Lock()
	defer m.lastImageMutex.Unlock()
	return m.lastImageBytes
}

func (m *MJPEGSplitter) setLastImageBytes(bytes []byte) {
	m.lastImageMutex.Lock()
	defer m.lastImageMutex.Unlock()
	m.lastImageBytes = bytes
}

func (m *MJPEGSplitter) inputHandlerLoop() {
	log.Print("Starting MJPEGSplitter inputHandlerLoop")
	buf := make([]byte, 65536)
	var streamBuffer []byte
	for !m.terminate {
		n, err := m.stream.Read(buf)
		log.Printf("Got %d bytes", n)
		if n > 0 {
			streamBuffer = append(streamBuffer, buf[:n]...)
			for {
				startIdx := bytes.Index(streamBuffer, soiMarker)
				if startIdx == -1 {
					break
				}
				endIdx := bytes.Index(streamBuffer[startIdx:], eoiMarker)
				if endIdx == -1 {
					break
				}
				endIdx += startIdx + len(eoiMarker)
				jpegBytes := streamBuffer[startIdx:endIdx]
				m.setLastImageBytes(jpegBytes)
				if m.listener != nil {
					m.listener(m)
				}
				streamBuffer = streamBuffer[endIdx:]
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
	}
	log.Print("Stopped MJPEGSplitter inputHandlerLoop")
}
