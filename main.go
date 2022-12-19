package main

import (
	"time"

	"github.com/anthonyme00/terminal-test/mandelbrot"
	"github.com/anthonyme00/terminal-test/output"
	"github.com/anthonyme00/terminal-test/window"
)

const (
	X_RES = 160
	Y_RES = 50

	FRAME_LIMIT = 60.0
)

var (
	startTime time.Time

	sleepDuration time.Duration
)

func init() {
	startTime = time.Now()

	sleepDurationInSecond := 1.0 / FRAME_LIMIT
	sleepDuration = time.Duration(float64(time.Second) * sleepDurationInSecond)
}

func main() {
	lastTime := time.Now()

	window := window.NewWindow(X_RES, Y_RES, '#', output.NewStdOutput())
	app := mandelbrot.NewMandelbrotApp()
	app.Init(window)

	for true {
		nextFrameTarget := lastTime.Add(sleepDuration)

		window.ClearScreen()
		app.Step(time.Millisecond * 1000)
		window.Draw()

		durationSleep := nextFrameTarget.Sub(time.Now())

		if durationSleep > 0 {
			time.Sleep(durationSleep - (time.Millisecond / 2))
		}

		for time.Now().Before(nextFrameTarget) {
		}

		lastTime = time.Now()
	}
}
