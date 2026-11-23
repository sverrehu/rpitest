module github.com/sverrehu/rpitest

go 1.27.0

replace github.com/sverrehu/rpigo => ./../rpigo

replace github.com/born-ml/born => ../../src/born.sverrehu

require (
	github.com/born-ml/born v0.9.24
	github.com/microsoft/onnxruntime/go v0.0.0-20261117191608-63163cb01727
	github.com/sverrehu/goutils v1.0.7
	github.com/sverrehu/rpigo v0.0.0
	github.com/yalue/onnxruntime_go v1.38.0
	gocv.io/x/gocv v0.43.0
	golang.org/x/image v0.48.0
)

require golang.org/x/sys v0.50.0 // indirect
