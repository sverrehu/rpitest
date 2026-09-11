package main

import (
	"log"
	"time"

	"github.com/sverrehu/rpigo/component"
	"github.com/sverrehu/rpigo/pwm"
)

const horizChan = 2
const vertChan = 3

func main() {
	gpio, err := component.NewGPIO()
	if err != nil {
		log.Panic(err)
	}
	defer gpio.Close()
	horizServo, err := newServo(horizChan)
	if err != nil {
		log.Panic(err)
	}
	defer horizServo.Close()
	vertServo, err := newServo(vertChan)
	if err != nil {
		log.Panic(err)
	}
	defer vertServo.Close()
	vertServo.SetAngle(0)
	time.Sleep(1 * time.Second)
	vertServo.SetAngle(45)
	time.Sleep(1 * time.Second)
	vertServo.SetAngle(90)
	time.Sleep(1 * time.Second)
	vertServo.SetAngle(135)
	time.Sleep(1 * time.Second)
	vertServo.SetAngle(180)
	time.Sleep(1 * time.Second)
	horizServo.SetAngle(90)
	vertServo.SetAngle(90)
	time.Sleep(1 * time.Second)
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
