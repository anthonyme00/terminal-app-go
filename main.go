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

	win := window.NewWindow(X_RES, Y_RES, '#', output.NewStdOutput())
	win.Open()
	defer win.Close()

	app := mandelbrot.NewMandelbrotApp()
	app.Init(win)

	for true {
		nextFrameTarget := lastTime.Add(sleepDuration)

		win.ClearScreen()
		app.Step(window.UpdatesInfo{
			Time_DeltaTime:    time.Now().Sub(lastTime),
			Time_AbsoluteTime: time.Now().Sub(startTime),
		})
		win.Draw()

		durationSleep := nextFrameTarget.Sub(time.Now())

		if durationSleep > 0 {
			time.Sleep(durationSleep - (time.Millisecond / 2))
		}

		for time.Now().Before(nextFrameTarget) {
		}

		lastTime = time.Now()
	}
}
