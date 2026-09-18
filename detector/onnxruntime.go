package detector

// https://github.com/yalue/onnxruntime_go_examples/blob/master/image_object_detect/image_object_detect.go

import (
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
	session   *ort.AdvancedSession
	input     *ort.Tensor[float32]
	output    *ort.Tensor[float32]
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
	err := ort.InitializeEnvironment()
	if err != nil {
		return err
	}
	inputShape := ort.NewShape(1, 3, int64(d.width), int64(d.height))
	inputTensor, err := ort.NewEmptyTensor[float32](inputShape)
	if err != nil {
		return err
	}
	outputShape := ort.NewShape(1, 300, 6)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		_ = inputTensor.Destroy()
		return err
	}
	options, err := ort.NewSessionOptions()
	if err != nil {
		_ = inputTensor.Destroy()
		_ = outputTensor.Destroy()
		return err
	}
	defer options.Destroy()

	session, err := ort.NewAdvancedSession(d.modelPath,
		[]string{"images"}, []string{"output0"},
		[]ort.ArbitraryTensor{inputTensor},
		[]ort.ArbitraryTensor{outputTensor},
		options)
	if err != nil {
		_ = inputTensor.Destroy()
		_ = outputTensor.Destroy()
		return err
	}
	d.session = session
	d.input = inputTensor
	d.output = outputTensor
	return nil
}

func (d *ONNXRuntimeDetector) Close() {
	if d.input != nil {
		_ = d.input.Destroy()
		d.input = nil
	}
	if d.output != nil {
		_ = d.output.Destroy()
		d.output = nil
	}
	if d.session != nil {
		_ = d.session.Destroy()
		d.session = nil
	}
}

func (d *ONNXRuntimeDetector) Detect(img *image.RGBA) ([]*Detection, error) {
	scale, err := d.loadAndProcessImage(img)
	if err != nil {
		return nil, err
	}
	st := time.Now()
	err = d.session.Run()
	if err != nil {
		return nil, err
	}
	log.Printf("session.Run time: %s\n", time.Since(st))
	detections := d.toDetections(scale)
	return detections, nil
}

func (d *ONNXRuntimeDetector) loadAndProcessImage(img *image.RGBA) (float32, error) {
	scaledImage, scale := scaleImage(img, d.width, d.height)
	data := d.input.GetData()
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

func (d *ONNXRuntimeDetector) toDetections(scale float32) []*Detection {
	// YOLO26 output format: [batch=1, num_detections=300, 6]
	// Each detection row: [x1, y1, x2, y2, score, class]
	data := d.output.GetData()
	numDetections := int(d.output.GetShape()[1])
	rowLen := int(d.output.GetShape()[2])
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
	return detections
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
