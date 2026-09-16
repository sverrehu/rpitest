package twoaxis

import (
	"fmt"
	"log"
	"time"

	"github.com/sverrehu/rpigo/component"
	"github.com/sverrehu/rpigo/pwm"
)

type TwoAxis struct {
	panChan      int
	tiltChan     int
	PanMinAngle  float64
	PanMaxAngle  float64
	TiltMinAngle float64
	TiltMaxAngle float64
	lastPan      float64
	lastTilt     float64
	terminate    bool
	panServo     *component.Servo
	tiltServo    *component.Servo
}

func NewTwoAxis(panChan, tiltChan int) *TwoAxis {
	return &TwoAxis{
		panChan:      panChan,
		tiltChan:     tiltChan,
		PanMinAngle:  0.0,
		PanMaxAngle:  180.0,
		TiltMinAngle: 65.0,
		TiltMaxAngle: 180.0,
		terminate:    false,
	}
}

func (ta *TwoAxis) Init() error {
	err := ta.setupServos()
	if err != nil {
		return err
	}
	return ta.Center()
}

func (ta *TwoAxis) Close() {
	ta.terminate = true
	ta.closeServos()
}

func (ta *TwoAxis) Pan(deg float64) error {
	if deg < ta.PanMinAngle || deg > ta.PanMaxAngle {
		return fmt.Errorf("pan %f not within valid range of [%f, %f]", deg, ta.PanMinAngle, ta.PanMaxAngle)
	}
	ta.lastPan = deg
	return ta.panServo.SetAngle(deg)
}

func (ta *TwoAxis) GetPan() float64 {
	return ta.lastPan
}

func (ta *TwoAxis) Tilt(deg float64) error {
	if deg < ta.TiltMinAngle || deg > ta.TiltMaxAngle {
		return fmt.Errorf("tilt %f not within valid range of [%f, %f]", deg, ta.PanMinAngle, ta.PanMaxAngle)
	}
	ta.lastTilt = deg
	return ta.tiltServo.SetAngle(deg)
}

func (ta *TwoAxis) GetTilt() float64 {
	return ta.lastTilt
}

func (ta *TwoAxis) PanTilt(p, t float64) error {
	err := ta.Pan(p)
	if err != nil {
		return err
	}
	return ta.Tilt(t)
}

func (ta *TwoAxis) Center() error {
	return ta.PanTilt(ta.PanMinAngle+(ta.PanMaxAngle-ta.PanMinAngle)/2.0, ta.TiltMinAngle+(ta.TiltMaxAngle-ta.TiltMinAngle)/2.0)
}

func (ta *TwoAxis) LimitPanTilt(pan, tilt float64) (float64, float64) {
	if pan < ta.PanMinAngle {
		pan = ta.PanMinAngle
	} else if pan > ta.PanMaxAngle {
		pan = ta.PanMaxAngle
	}
	if tilt < ta.TiltMinAngle {
		tilt = ta.TiltMinAngle
	} else if tilt > ta.TiltMaxAngle {
		tilt = ta.TiltMaxAngle
	}
	return pan, tilt
}

func (ta *TwoAxis) newServo(channel int) (*component.Servo, error) {
	p, err := pwm.NewHardPWM(0, channel, 50)
	if err != nil {
		return nil, err
	}
	servo, err := component.NewServo(p)
	if err != nil {
		return nil, err
	}
	return servo, nil
}

func (ta *TwoAxis) setupServos() error {
	var err error
	ta.panServo, err = ta.newServo(ta.panChan)
	if err != nil {
		return err
	}
	ta.tiltServo, err = ta.newServo(ta.tiltChan)
	if err != nil {
		return err
	}
	return nil
}

func (ta *TwoAxis) closeServos() {
	err := ta.Center()
	if err != nil {
		log.Printf("Error centering: %v", err)
	}
	time.Sleep(500 * time.Millisecond)
	if ta.tiltServo != nil {
		_ = ta.tiltServo.Close()
		ta.tiltServo = nil
	}
	if ta.panServo != nil {
		_ = ta.panServo.Close()
		ta.panServo = nil
	}
}
