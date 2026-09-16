package main

import (
	"context"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sverrehu/rpitest/camera"
	"github.com/sverrehu/rpitest/detector"
	"github.com/sverrehu/rpitest/imgdraw"
	"github.com/sverrehu/rpitest/imgposter"
	"github.com/sverrehu/rpitest/twoaxis"
)

const taPanChan = 2
const taTiltChan = 3
const panDegreesPerImage = 70.0
const tiltDegreesPerImage = 40.0

var ta *twoaxis.TwoAxis

func main() {
	rpiviewHost := "192.168.1.15"
	//rpiviewHost := "192.168.30.21"
	rpiview := imgposter.NewImagePoster(rpiviewHost, 8086)

	ta = twoaxis.NewTwoAxis(taPanChan, taTiltChan)
	err := ta.Init()
	if err != nil {
		log.Panic(err)
	}
	installTerminationHandler2()

	det := detector.NewDetector("../gotest/gocv/yolo26_face_fp16.onnx", 640, 640)
	err = det.Init()
	if err != nil {
		panic(err)
	}
	cam := camera.NewCamera()
	defer cam.Close()
	cam.Rotate = true
	err = cam.StartStreaming(nil)
	if err != nil {
		panic(err)
	}
	for {
		img, err := cam.GetImage()
		if img == nil {
			continue
		}
		if err != nil {
			log.Panic(err)
		}
		detections, err := det.Detect(img)
		if err != nil {
			panic(err)
		}
		annotateDetections2(img, detections)
		err = rpiview.PostJPEG(img)
		if err != nil {
			log.Panic(err)
		}
		detection := selectDetection(detections)
		if detection != nil {
			err := moveTo(detection, img)
			if err != nil {
				panic(err)
			}
		} else {
			err := ta.Center()
			if err != nil {
				panic(err)
			}
		}
		// give servos time to settle before grabbing the next image
		time.Sleep(300 * time.Millisecond)
	}
}

func moveTo(detection *detector.Detection, img *image.RGBA) error {
	if detection.X1 < detection.X0 || detection.Y1 < detection.Y0 {
		panic("Wrong order!")
	}
	bounds := img.Bounds()
	panDegreesPerPixel := panDegreesPerImage / float64(bounds.Dx())
	tiltDegreesPerPixel := tiltDegreesPerImage / float64(bounds.Dy())
	roiHitRadius := math.Min(float64(bounds.Dx()), float64(bounds.Dy())) * 0.02
	imageCenterX := bounds.Dx() / 2
	imageCenterY := bounds.Dy() / 2
	roiCenterX := detection.X0 + (detection.X1-detection.X0)/2
	roiCenterY := detection.Y0 + (detection.Y1-detection.Y0)/2
	dx := imageCenterX - roiCenterX
	dy := imageCenterY - roiCenterY
	if math.Abs(float64(dx)) < roiHitRadius {
		dx = 0
	}
	if math.Abs(float64(dy)) < roiHitRadius {
		dy = 0
	}
	pan := ta.GetPan()
	tilt := ta.GetTilt()
	wantedPan := pan - float64(-dx)*panDegreesPerPixel
	wantedTilt := tilt - float64(dy)*tiltDegreesPerPixel
	wantedPan, wantedTilt = ta.LimitPanTilt(wantedPan, wantedTilt)
	log.Printf("dx: %d, dy: %d, pan: %f -> %f, tilt: %f -> %f", dx, dy, pan, wantedPan, tilt, wantedTilt)
	return ta.PanTilt(wantedPan, wantedTilt)
}

func selectDetection(detections []*detector.Detection) *detector.Detection {
	biggestArea := -1.0
	var biggest *detector.Detection
	for _, d := range detections {
		w := math.Abs(float64(d.X1 - d.X0))
		h := math.Abs(float64(d.Y1 - d.Y0))
		area := w * h
		if area > biggestArea {
			biggestArea = area
			biggest = d
		}
	}
	return biggest
}

func annotateDetections2(img image.Image, detections []*detector.Detection) {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		panic("Image is not RGBA")
	}
	id := imgdraw.NewImageDrawer(rgba)
	bounds := rgba.Bounds()
	id.Color(color.RGBA{255, 100, 100, 255})
	id.Line(0, bounds.Dy()/2, bounds.Dx(), bounds.Dy()/2)
	id.Line(bounds.Dx()/2, 0, bounds.Dx()/2, bounds.Dy())
	id.Color(color.RGBA{100, 255, 100, 255})
	for _, d := range detections {
		id.Rect(d.X0, d.Y0, d.X1, d.Y1)
	}
}

func installTerminationHandler2() {
	go func() {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
		<-ctx.Done()
		log.Println("Shutting down...")
		ta.Close()
		log.Println("Done.")
		os.Exit(0)
	}()
}
