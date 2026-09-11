package main

import (
	"log"
	"time"

	"github.com/sverrehu/rpigo/component"
	"github.com/sverrehu/rpigo/pwm"
)

const panChan = 2
const tiltChan = 3
const panMinAngle = 0
const panMaxAngle = 180
const tiltMinAngle = 65 // points up
const tiltMaxAngle = 180

var panServo *component.Servo
var tiltServo *component.Servo

func main() {
	err := setupServos()
	if err != nil {
		log.Panic(err)
	}
	defer close()
	dPan := 0.2
	dTilt := 0.2
	pan := panMinAngle + (panMaxAngle-panMinAngle)/2.0
	tilt := tiltMinAngle + (tiltMaxAngle-tiltMinAngle)/2.0
	for {
		log.Printf("Pan: %f, Tilt: %f, dPan: %f, dTilt: %f", pan, tilt, dPan, dTilt)
		err := panServo.SetAngle(pan)
		if err != nil {
			log.Panic(err)
		}
		err = tiltServo.SetAngle(tilt)
		if err != nil {
			log.Panic(err)
		}
		pan += dPan
		if pan < panMinAngle || pan > panMaxAngle {
			dPan = -dPan
			pan += dPan
		}
		tilt += dTilt
		if tilt < tiltMinAngle || tilt > tiltMaxAngle {
			dTilt = -dTilt
			tilt += dTilt
		}
		time.Sleep(1 * time.Millisecond)
	}
	center()
	log.Println("Servos reset. Sleeping a little.")
	time.Sleep(5 * time.Second)
}

func newServo(channel int) (*component.Servo, error) {
	pwm, err := pwm.NewHardPWM(0, channel, 50)
	if err != nil {
		return nil, err
	}
	servo, err := component.NewServo(pwm)
	if err != nil {
		return nil, err
	}
	return servo, nil
}

func setupServos() error {
	var err error
	panServo, err = newServo(panChan)
	if err != nil {
		return err
	}
	tiltServo, err = newServo(tiltChan)
	if err != nil {
		return err
	}
	return nil
}

func close() {
	if tiltServo != nil {
		_ = tiltServo.Close()
	}
	if panServo != nil {
		_ = panServo.Close()
	}
}

func center() error {
	err := panServo.SetAngle(panMinAngle + (panMaxAngle-panMinAngle)/2)
	if err != nil {
		return err
	}
	err = tiltServo.SetAngle(tiltMinAngle + (tiltMaxAngle-tiltMinAngle)/2)
	if err != nil {
		return err
	}
	return nil
}
