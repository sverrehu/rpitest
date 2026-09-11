package main

import (
	"log"
	"time"

	"github.com/sverrehu/rpigo/component"
	"github.com/sverrehu/rpigo/pwm"
)

func main() {
	gpio, err := component.NewGPIO()
	if err != nil {
		log.Panic(err)
	}
	defer gpio.Close()
	pwm, err := pwm.NewHardPWM(0, 2, 50)
	if err != nil {
		log.Panic(err)
	}
	servo, err := component.NewServo(pwm)
	if err != nil {
		log.Panic(err)
	}
	defer servo.Close()
	servo.SetAngle(0)
	time.Sleep(1 * time.Second)
	servo.SetAngle(45)
	time.Sleep(1 * time.Second)
	servo.SetAngle(90)
	time.Sleep(1 * time.Second)
	servo.SetAngle(135)
	time.Sleep(1 * time.Second)
	servo.SetAngle(180)
	time.Sleep(1 * time.Second)
}
