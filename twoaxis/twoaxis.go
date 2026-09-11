package main

import (
	"log"
	"time"

	"github.com/sverrehu/rpigo/component"
	"github.com/sverrehu/rpigo/pwm"
)

const horizChan = 2
const vertChan = 3
const horizMinAngle = 0
const horizMaxAngle = 180
const vertMinAngle = 65 // points up
const vertMaxAngle = 180

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
	for angle := 0; angle <= 180; angle++ {
		horizServo.SetAngle(float64(angle))
		time.Sleep(100 * time.Millisecond)
	}
	horizServo.SetAngle((horizMaxAngle - horizMinAngle) / 2)
	vertServo.SetAngle((vertMaxAngle - vertMinAngle) / 2)
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
