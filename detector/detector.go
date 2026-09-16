package detector

import (
	"image"

	"github.com/born-ml/born/backend/cpu"
	"github.com/born-ml/born/onnx"
	"github.com/born-ml/born/tensor"
)

type Detector struct {
	modelPath string
	width     int
	height    int
	backend   tensor.Backend
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
	d.backend = be
	d.model = model
	return nil
}

func (d *Detector) Close() {
}

func (d *Detector) Detect(img *image.RGBA) ([]Detection, error) {
	it, err := d.loadAndProcessImage(img)
	if err != nil {
		return nil, err
	}
	_, err = d.model.Forward(it.Raw())
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (d *Detector) loadAndProcessImage(img *image.RGBA) (*tensor.Tensor[float32, tensor.Backend], error) {
	data := make([]float32, 3*d.width*d.height)
	offsetX := 0
	offsetY := 0
	for y := 0; y < d.height; y++ {
		for x := 0; x < d.width; x++ {
			r, g, b, _ := img.At(offsetX+x, offsetY+y).RGBA()

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
		return nil, err
	}
	return t, nil
}
