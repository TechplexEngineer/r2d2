package main

import (
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/joystick"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	throttleZero = 350
)

// adjusts the throttle from -1.0 (hard back) <-> 1.0 (hard forward) to the correct
// pwm pulse values.
func throttleToPwmWidth(val float64) int {
	if val > 0 {
		return int(gobot.Rescale(val, 0, 1, 350, 300))
	}
	return int(gobot.Rescale(val, -1, 0, 490, 350))
}

func main() {
	rpi := raspi.NewAdaptor()
	//led := gpio.NewLedDriver(rpi, "7")
	joystickAdaptor := joystick.NewAdaptor()
	stick := joystick.NewDriver(joystickAdaptor, joystick.Xbox360)
	pwm := i2c.NewPCA9685Driver(rpi,
		i2c.WithBus(0),
		i2c.WithAddress(0x34))

	leftSpeed := 0.0
	rightSpeed := 0.0

	work := func() {
		stick.On(joystick.LeftX, func(data interface{}) {
			leftSpeed = -float64(data.(int16)) / 10000
		})
		stick.On(joystick.RightY, func(data interface{}) {
			rightSpeed = -float64(data.(int16)) / 10000
		})

		// init the PWM controller
		pwm.SetPWMFreq(60)

		// init the ESC controller for throttle zero
		pwm.SetPWM(0, 0, uint16(throttleZero))
		pwm.SetPWM(1, 0, uint16(throttleZero))

		gobot.Every(20*time.Millisecond, func() {

			pwm.SetPWM(0, 0, uint16(throttleToPwmWidth(leftSpeed)))
			pwm.SetPWM(1, 0, uint16(throttleToPwmWidth(rightSpeed)))

		})
		//gobot.Every(1*time.Second, func() {
		//	//led.Toggle()
		//	//stick.Subscribe()
		//	log.Print("Hello, other")
		//})
	}

	robot := gobot.NewRobot("r2d2",
		[]gobot.Connection{rpi},
		[]gobot.Device{
			stick,
			pwm,
		},
		work,
	)

	robot.Start()
}
