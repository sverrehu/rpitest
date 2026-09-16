package detector

import (
	"image"
	"math"

	"github.com/born-ml/born/backend/cpu"
	"github.com/born-ml/born/onnx"
	"github.com/born-ml/born/tensor"
	"golang.org/x/image/draw"
)

type Detector struct {
	modelPath string
	width     int
	height    int
	backend   tensor.Backend
	model     onnx.Model
}

type Detection struct {
	X0, Y0 int
	X1, Y1 int
	Class  int
	Score  float32
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

func (d *Detector) Detect(img *image.RGBA) ([]*Detection, error) {
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

func (d *Detector) loadAndProcessImage(img *image.RGBA) (*tensor.Tensor[float32, tensor.Backend], float32, error) {
	scaledImage, scale := d.scaleImage(img)
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

func (d *Detector) scaleImage(img *image.RGBA) (*image.RGBA, float32) {
	dst := image.NewRGBA(image.Rect(0, 0, d.width, d.height))
	hScale := float64(img.Bounds().Dx()) / float64(d.width)
	vScale := float64(img.Bounds().Dy()) / float64(d.height)
	scale := float32(math.Max(hScale, vScale))
	dstRect := image.Rect(0, 0, int(float32(img.Bounds().Dx())/scale), int(float32(img.Bounds().Dy())/scale))
	draw.ApproxBiLinear.Scale(dst, dstRect, img, img.Bounds(), draw.Over, nil)
	return dst, scale
}

func (d *Detector) toDetections(t *tensor.RawTensor, scale float32) []*Detection {
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
