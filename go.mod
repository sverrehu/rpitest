module github.com/sverrehu/rpitest

go 1.27.0

replace github.com/sverrehu/rpigo => ./../rpigo

replace github.com/born-ml/born => ../../src/born.sverrehu

require (
	github.com/born-ml/born v0.9.23
	github.com/sverrehu/rpigo v0.0.0
	gocv.io/x/gocv v0.43.0
	golang.org/x/image v0.46.0
	github.com/yalue/onnxruntime_go v1.36.0 // indirect
)

require golang.org/x/sys v0.48.0 // indirect
