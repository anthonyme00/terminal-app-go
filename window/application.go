package window

import (
	"time"

	"github.com/anthonyme00/terminal-test/output"
)

type Application interface {
	Init(windowData *Window)
	Step(delta time.Duration)
}

type Window struct {
	screen     []byte
	xSize      int
	ySize      int
	BorderChar byte
	Output     output.Output
}

func NewWindow(x_size, y_size int, border byte, output output.Output) *Window {
	// Add 3 to x size for window border + new line character
	// Add 2 to y size for window border
	screen := make([]byte, (x_size+3)*(y_size+2))

	return &Window{
		xSize:      x_size,
		ySize:      y_size,
		BorderChar: border,
		screen:     screen,
		Output:     output,
	}
}

func (w *Window) coordToIndex(x_abs, y_abs int) int {
	if x_abs >= w.xSize+3 {
		x_abs = w.xSize + 2
	}

	if y_abs >= w.ySize+2 {
		y_abs = w.ySize + 1
	}

	return y_abs*(w.xSize+3) + x_abs
}

func (w *Window) SetScreen(x, y int, char byte) {
	w.screen[w.coordToIndex(x+1, y+1)] = char
}

func (w *Window) Dump() []byte {
	return w.screen
}

func (w *Window) GetSize() (x, y int) {
	return w.xSize, w.ySize
}

func (w *Window) ClearScreen() {
	// draw default character
	for i := 0; i < len(w.screen); i++ {
		w.screen[i] = w.BorderChar
	}

	// add newlines
	for y := 0; y < w.ySize+2; y++ {
		w.screen[w.coordToIndex(w.xSize+2, y)] = '\n'
	}

	// empties drawing canvas
	for x := 0; x < w.xSize; x++ {
		for y := 0; y < w.ySize; y++ {
			w.screen[w.coordToIndex(x+1, y+1)] = ' '
		}
	}
}

func (w *Window) Draw() {
	w.Output.Write(w.screen)
}
