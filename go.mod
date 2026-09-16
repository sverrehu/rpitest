module github.com/sverrehu/rpitest

go 1.27.0

replace github.com/born-ml/born => ../../src/born.sverrehu

replace github.com/sverrehu/rpitest/camera => ./camera

require (
	github.com/born-ml/born v0.9.23
	github.com/sverrehu/rpitest/camera v0.0.0
	gocv.io/x/gocv v0.43.0
	golang.org/x/image v0.46.0
)

require golang.org/x/sys v0.48.0 // indirect
