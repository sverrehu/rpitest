package detector

import (
	"image"

	"github.com/born-ml/born/backend/cpu"
	"github.com/born-ml/born/onnx"
)

type Detector struct {
	modelPath string
	width     int
	height    int
	model     onnx.Model
}

type Detection struct {
	x0, y0 int
	x1, y1 int
}

func NewDetector(modelPath string, width int, height int) *Detector {
	return &Detector{
		modelPath: modelPath,
		width:     width,
		height:    height,
	}
}

func (d *Detector) Init() error {
	be := cpu.New()
	model, err := onnx.Load(d.modelPath, be)
	if err != nil {
		return err
	}
	d.model = model
	return nil
}

func (d *Detector) Close() {
}

func (d *Detector) Detect(img *image.RGBA) ([]Detection, error) {
	return nil, nil
}
