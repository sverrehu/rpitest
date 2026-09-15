module github.com/sverrehu/rpitest

go 1.27.0

replace github.com/sverrehu/rpitest/camera => ./camera

require (
	github.com/sverrehu/rpitest/camera v0.0.0
	gocv.io/x/gocv v0.43.0
)
