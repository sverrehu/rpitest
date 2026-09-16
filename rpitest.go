package main

// On my Mac: PKG_CONFIG_PATH=/opt/local/lib/opencv4/pkgconfig go run viewvideo.go

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"gocv.io/x/gocv"
)

func main() {
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Fatalf("Error opening web cam: %v", err)
	}
	defer webcam.Close()

	modelPath := "yolo26_face_fp16.onnx"
	net := gocv.ReadNetFromONNX(modelPath)
	if net.Empty() {
		log.Fatalf("Error reading network from model file: %s", modelPath)
	}
	defer net.Close()

	net.SetPreferableBackend(gocv.NetBackendDefault)
	net.SetPreferableTarget(gocv.NetTargetCPU)

	window := gocv.NewWindow("GoCV YOLO26 Mac Camera")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	fmt.Println("Press 'q' in the camera window to exit.")

	for {
		if ok := webcam.Read(&img); !ok || img.Empty() {
			log.Println("Device closed or unable to read from the webcam")
			break
		}

		blob := gocv.BlobFromImage(img, 1.0/255.0, image.Pt(640, 640), gocv.NewScalar(0, 0, 0, 0), true, false)
		net.SetInput(blob, "")

		outputs := net.Forward("")

		processDetections(&img, outputs)

		blob.Close()
		outputs.Close()

		window.IMShow(img)
		if window.WaitKey(1) == int('q') {
			break
		}
	}
}

func processDetections(frame *gocv.Mat, outputs gocv.Mat) {
	confidenceThreshold := float32(0.25)
	green := color.RGBA{0, 255, 0, 0}

	for i := 0; i < outputs.Rows(); i++ {
		confidence := outputs.GetFloatAt(i, 4)
		if confidence > confidenceThreshold {
			centerX := outputs.GetFloatAt(i, 0)
			centerY := outputs.GetFloatAt(i, 1)
			width := outputs.GetFloatAt(i, 2)
			height := outputs.GetFloatAt(i, 3)

			scaleX := float32(frame.Cols()) / 640.0
			scaleY := float32(frame.Rows()) / 640.0

			left := int((centerX - width/2) * scaleX)
			top := int((centerY - height/2) * scaleY)
			right := int((centerX + width/2) * scaleX)
			bottom := int((centerY + height/2) * scaleY)

			rect := image.Rect(left, top, right, bottom)
			gocv.Rectangle(frame, rect, green, 2)

			gocv.PutText(frame, fmt.Sprintf("Object: %.2f", confidence), image.Pt(left, top-10),
				gocv.FontHersheySimplex, 0.5, green, 2)
		}
	}
}
