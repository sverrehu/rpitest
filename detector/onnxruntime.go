package detector

// https://github.com/yalue/onnxruntime_go_examples/blob/master/image_object_detect/image_object_detect.go

import (
	"context"
	"fmt"
	"image"
	"log"
	"runtime"
	"time"

	ort "github.com/microsoft/onnxruntime/go/onnxruntime"
)

const onnxruntimeVersion = "1.30.0"
const onnxruntimeLibPath = "../../lib/onnxruntime"

type ONNXRuntimeDetector struct {
	modelPath string
	width     int
	height    int
	session   *ort.Session
	input     *ort.Tensor
	inputData []float32
}

func NewONNXRuntimeDetector(modelPath string, width int, height int) *ONNXRuntimeDetector {
	return &ONNXRuntimeDetector{
		modelPath: modelPath,
		width:     width,
		height:    height,
	}
}

func (d *ONNXRuntimeDetector) Init() error {
	ort.SetSharedLibraryPath(findSharedLibrary())
	err := ort.Init()
	if err != nil {
		return err
	}
	inputShape := []int64{1, 3, int64(d.width), int64(d.height)}
	d.inputData = make([]float32, 3*d.width*d.height)
	d.input, err = ort.CreateTensor[float32](inputShape, d.inputData)
	if err != nil {
		return err
	}
	options, err := ort.NewSessionOptions()
	if err != nil {
		d.Close()
		return err
	}
	defer options.Close()
	err = options.AppendExecutionProvider("WebGpuExecutionProvider", nil)
	if err != nil {
		d.Close()
		return err
	}

	d.session, err = ort.NewSession(d.modelPath, options)
	if err != nil {
		d.Close()
		return err
	}
	return nil
}

func (d *ONNXRuntimeDetector) Close() {
	if d.input != nil {
		_ = d.input.Close()
		d.input = nil
	}
	if d.session != nil {
		_ = d.session.Close()
		d.session = nil
	}
	_ = ort.Shutdown()
}

func (d *ONNXRuntimeDetector) Detect(img *image.RGBA) ([]*Detection, error) {
	scale, err := d.loadAndProcessImage(img)
	if err != nil {
		return nil, err
	}
	st := time.Now()
	results, err := d.session.Run(context.Background(), map[string]*ort.Tensor{
		"images": d.input,
	}, nil)
	if err != nil {
		return nil, err
	}
	outTensor := results["output0"]
	if outTensor == nil {
		return nil, fmt.Errorf("output0 is nil")
	}
	defer outTensor.Close()
	log.Printf("session.Run time: %s\n", time.Since(st))
	return d.toDetections(outTensor, scale)
}

func (d *ONNXRuntimeDetector) loadAndProcessImage(img *image.RGBA) (float32, error) {
	scaledImage, scale := scaleImage(img, d.width, d.height)
	data := d.inputData
	offsetX := 0
	offsetY := 0
	for y := 0; y < d.height; y++ {
		for x := 0; x < d.width; x++ {
			r, g, b, _ := scaledImage.At(offsetX+x, offsetY+y).RGBA()
			idxR := 0*d.width*d.height + y*d.height + x
			idxG := 1*d.width*d.height + y*d.height + x
			idxB := 2*d.width*d.height + y*d.height + x
			data[idxR] = float32(r) / 65535.0
			data[idxG] = float32(g) / 65535.0
			data[idxB] = float32(b) / 65535.0
		}
	}
	return scale, nil
}

func (d *ONNXRuntimeDetector) toDetections(outTensor *ort.Tensor, scale float32) ([]*Detection, error) {
	// YOLO26 output format: [batch=1, num_detections=300, 6]
	// Each detection row: [x1, y1, x2, y2, score, class]
	data, err := ort.TensorData[float32](outTensor)
	if err != nil {
		return nil, err
	}
	numDetections := int(outTensor.Shape()[1])
	rowLen := int(outTensor.Shape()[2])
	confidenceThreshold := float32(0.25)
	detections := make([]*Detection, 0)
	for i := 0; i < numDetections; i++ {
		offset := i * rowLen
		x0 := data[offset+0]
		y0 := data[offset+1]
		x1 := data[offset+2]
		y1 := data[offset+3]
		score := data[offset+4]
		clss := int(data[offset+5])
		if score >= confidenceThreshold {
			d := &Detection{int(x0 * scale), int(y0 * scale), int(x1 * scale), int(y1 * scale), clss, score}
			detections = append(detections, d)
		}
	}
	return detections, nil
}

func findSharedLibrary() string {
	if runtime.GOOS == "darwin" {
		if runtime.GOARCH == "arm64" {
			return fmt.Sprintf("%s/onnxruntime-osx-arm64-%s/lib/libonnxruntime.%s.dylib", onnxruntimeLibPath, onnxruntimeVersion, onnxruntimeVersion)
		}
	}
	if runtime.GOOS == "linux" {
		if runtime.GOARCH == "arm64" {
			return fmt.Sprintf("%s/onnxruntime-linux-aarch64-%s/lib/libonnxruntime.so.%s", onnxruntimeLibPath, onnxruntimeVersion, onnxruntimeVersion)
		}
	}
	log.Fatalf("Unable to determine a path to the onnxruntime shared library for OS \"%s\" and architecture \"%s\".\n",
		runtime.GOOS, runtime.GOARCH)
	return "//never gets here"
}
