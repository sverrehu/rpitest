package detector

import (
	"fmt"
	"image"
	"log"
	"runtime"

	"github.com/born-ml/born/onnx"
	"github.com/born-ml/born/tensor"
	ort "github.com/yalue/onnxruntime_go"
)

const onnxruntimeVersion = "1.29.0"
const onnxruntimeLibPath = "../../../lib/onnxruntime"

type ONNXRuntimeDetector struct {
	modelPath string
	width     int
	height    int
	backend   tensor.Backend
	model     onnx.Model
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
	return nil
}

func (d *ONNXRuntimeDetector) Close() {
	_ = ort.DestroyEnvironment()
}

func (d *ONNXRuntimeDetector) Detect(img *image.RGBA) ([]*Detection, error) {
	it, scale, err := d.loadAndProcessImage(img)
	if err != nil {
		return nil, err
	}
	ot, err := d.model.Forward(it.Raw())
	if err != nil {
		return nil, err
	}
	detections := d.toDetections(ot, scale)
	return detections, nil
}

func (d *ONNXRuntimeDetector) loadAndProcessImage(img *image.RGBA) (*tensor.Tensor[float32, tensor.Backend], float32, error) {
	scaledImage, scale := scaleImage(img, d.width, d.height)
	data := make([]float32, 3*d.width*d.height)
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
	t, err := tensor.FromSlice(data, tensor.Shape{1, 3, 640, 640}, d.backend)
	if err != nil {
		return nil, 0, err
	}
	return t, scale, nil
}

func (d *ONNXRuntimeDetector) toDetections(t *tensor.RawTensor, scale float32) []*Detection {
	// YOLO26 output format: [batch=1, num_detections=300, 6]
	// Each detection row: [x1, y1, x2, y2, score, class]
	data := t.AsFloat32()
	numDetections := t.Shape()[1]
	rowLen := t.Shape()[2]
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
