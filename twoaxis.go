package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sverrehu/rpitest/twoaxis"
)

const panChan = 2
const tiltChan = 3

var axis *twoaxis.TwoAxis

func main() {
	axis = twoaxis.NewTwoAxis(panChan, tiltChan)
	err := axis.Init()
	if err != nil {
		log.Panic(err)
	}
	installTerminationHandler()
	dPan := 0.3
	dTilt := 0.3
	pan := axis.GetPan()
	tilt := axis.GetTilt()
	for {
		err := axis.PanTilt(pan, tilt)
		if err != nil {
			log.Panic(err)
		}
		pan += dPan
		if pan < axis.PanMinAngle || pan > axis.PanMaxAngle {
			dPan = -dPan
			pan += dPan
		}
		tilt += dTilt
		if tilt < axis.TiltMinAngle || tilt > axis.TiltMaxAngle {
			dTilt = -dTilt
			tilt += dTilt
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func installTerminationHandler() {
	go func() {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
		<-ctx.Done()
		log.Println("Shutting down...")
		axis.Close()
		log.Println("Done.")
		os.Exit(0)
	}()
}
