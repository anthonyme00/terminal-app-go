package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/anthonyme00/terminal-test/mandelbrot"
	tty "github.com/mattn/go-tty"
)

const (
	DEG_2_RAD = math.Pi / 180.0
	X_RES     = 160
	Y_RES     = 50

	FRAME_CHAR = byte('*')

	PRINT_CURSOR_POS = "\033[6n"

	CATCH_MOUSE_EVENTS   = "\033[?1003h\033[?1015h\033[?1006h"
	DISABLE_MOUSE_EVENTS = "\033[?1000l"
)

var (
	SCREENBUF  *[X_RES + 3][Y_RES + 2]byte = &[X_RES + 3][Y_RES + 2]byte{}
	TIME_START time.Time
	BUILDER    = strings.Builder{}

	FRAME_LIMIT    = 60.0
	SLEEP_DURATION time.Duration
)

func init() {
	clear()
	TIME_START = time.Now()
	BUILDER.Grow((X_RES + 3) * (Y_RES + 2))

	sleep_secs := 1.0 / FRAME_LIMIT
	SLEEP_DURATION = time.Duration(float64(time.Second) * sleep_secs)
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func draw() {
	BUILDER.Reset()

	byteBuf := make([]byte, X_RES+3)
	for y := 0; y < Y_RES+2; y++ {
		for x := 0; x < X_RES+3; x++ {
			byteBuf[x] = SCREENBUF[x][y]
		}

		BUILDER.Write(byteBuf)
	}

	fmt.Print(BUILDER.String())
}

func init_screen() {
	for x := 0; x < X_RES+3; x++ {
		for y := 0; y < Y_RES+2; y++ {
			SCREENBUF[x][y] = byte(' ')

			if x == 0 || y == 0 || x == X_RES+1 || y == Y_RES+1 {
				SCREENBUF[x][y] = FRAME_CHAR
			}

			if x == X_RES+2 {
				SCREENBUF[x][y] = byte('\n')
			}
		}
	}
}

func drawToScreen(x, y int, char byte) {
	SCREENBUF[x+1][y+1] = char
}

func write(x, y int, text string) {
	l_x := x
	l_y := y

	for i := 0; i < len(text); i++ {
		if l_x >= X_RES {
			l_x = x
			l_y += 1
		}

		if l_y >= Y_RES {
			return
		}

		drawToScreen(l_x, l_y, text[i])

		l_x += 1
	}
}

func step_wave() {
	type WaveDefFill struct {
		Fill bool
		Char byte
	}
	type WaveDef struct {
		Timescale float64
		Period    float64
		Offset    float64
		Amplitude float64
		Char      byte
		Fill      WaveDefFill
	}

	waves := []WaveDef{
		{
			Timescale: 1.7,
			Period:    0.213,
			Offset:    0.2,
			Char:      byte('$'),
			Amplitude: 1.0,
			Fill: WaveDefFill{
				Fill: true,
				Char: byte('@'),
			},
		},
		{
			Timescale: 1.5,
			Period:    1.333,
			Offset:    0.5,
			Char:      byte('^'),
			Amplitude: 0.7,
			Fill: WaveDefFill{
				Fill: true,
				Char: byte(';'),
			},
		},
		{
			Timescale: 1.0,
			Period:    1.0,
			Offset:    0,
			Char:      byte('+'),
			Amplitude: 0.4,
			Fill: WaveDefFill{
				Fill: true,
				Char: ('.'),
			},
		},
	}

	currentTime := time.Now()
	absoluteStartTime := currentTime.Sub(TIME_START).Seconds()

	for _, wave := range waves {
		for x := 0; x < X_RES; x++ {
			relativePos := float64(x-1) / float64(X_RES-1)
			relativeTime := relativePos * wave.Timescale

			absoluteTime := absoluteStartTime + relativeTime + wave.Offset

			yPos := math.Sin(2 * math.Pi * absoluteTime * wave.Period)
			yPos = (yPos + 1.0) / 2.0
			yPos = yPos * wave.Amplitude

			y := int((1-wave.Amplitude)*(Y_RES-1)) + int(math.Round(yPos*float64(Y_RES-1)))

			if wave.Fill.Fill {
				for j := y + 1; j < Y_RES; j++ {
					drawToScreen(x, j, wave.Fill.Char)
				}
			}

			drawToScreen(x, int(y), wave.Char)
		}
	}
}

func step_mandelbrot() {
	mandelbrotPeriod := 30.0

	lerp := func(a, b, t float64) float64 {
		return a + ((b - a) * t)
	}

	rotate := func(x, y, deltaDeg float64) (float64, float64) {
		deltaRad := deltaDeg * DEG_2_RAD
		return x*math.Cos(deltaRad) - y*math.Sin(deltaRad), x*math.Sin(deltaRad) + y*math.Cos(deltaRad)
	}

	rotateAround := func(o_x, o_y, x, y, deltaDeg float64) (float64, float64) {
		x_ref := x - o_x
		y_ref := y - o_y

		x_n, y_n := rotate(x_ref, y_ref, deltaDeg)

		return o_x + x_n, o_y + y_n
	}

	terminal_font_y_correction_factor := 2.0

	graphic := "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
	center_x := -0.7457
	center_y := 0.1127

	x_width_min := 2.00000
	x_width_max := 0.00001

	minRot := 0
	maxRot := 360

	currentTime := time.Now()
	absoluteTime := currentTime.Sub(TIME_START).Seconds()

	t := math.Pow(math.Sin(math.Pi*(absoluteTime-mandelbrotPeriod/2.0)/mandelbrotPeriod), 4)

	aspect_ratio := float64(X_RES) / float64(Y_RES)
	x_width := lerp(x_width_max, x_width_min, t)
	y_width := x_width / aspect_ratio * terminal_font_y_correction_factor

	getCoord := func(x, y int) (float64, float64) {
		rel_x := float64(x) / float64(X_RES-1)
		rel_y := float64(y) / float64(Y_RES-1)

		x_start, x_end := center_x-x_width, center_x+x_width
		y_start, y_end := center_y-y_width, center_y+y_width

		return lerp(x_start, x_end, rel_x), lerp(y_start, y_end, rel_y)
	}

	for x := 0; x < X_RES; x++ {
		for y := 0; y < Y_RES; y++ {
			s_x, s_y := getCoord(x, y)

			s_x, s_y = rotateAround(center_x, center_y, s_x, s_y, lerp(float64(minRot), float64(maxRot), t))

			m := mandelbrot.Sample(s_x, s_y, 100)
			i := int(m * (float64(len(graphic)-1) + 0.1))
			drawToScreen(x, y, graphic[i])
		}
	}
}

func drawCircle(x_p, y_p int, radius, ringRad float64) {
	terminal_font_y_correction_factor := 2.0

	graphic := ";:,\"^`'. "
	rSqr := radius * radius

	waveFreq := 3.0
	waveSpeed := 1.0

	absoluteTime := time.Now().Sub(TIME_START).Seconds()

	t := math.Mod(absoluteTime, waveFreq)

	for x := 0; x < X_RES; x++ {
		for y := 0; y < Y_RES; y++ {
			x_abs := float64(x_p - x)
			y_abs := float64(y_p - y)

			distSqr := (x_abs * x_abs) + ((y_abs * y_abs) * (terminal_font_y_correction_factor * terminal_font_y_correction_factor))

			targetDist := t * rSqr * waveSpeed

			if distSqr <= rSqr && math.Abs(distSqr-targetDist) < ringRad {
				i := float64(len(graphic)-1) * (math.Abs(distSqr-targetDist) / ringRad)
				drawToScreen(x, y, graphic[int(i)])
			}
		}
	}
}

func main() {
	term, _ := tty.Open()
	defer term.Close()

	counter := 0
	frameTime := [20]float64{}

	baseTime := time.Now()
	lastTime := time.Now()

	lastEvent := ""
	mouseX := 0
	mouseY := 0

	for true {
		init_screen()
		step_wave()
		// step_mandelbrot()

		sumFrameTime := 0.0
		for i := 0; i < 20; i++ {
			sumFrameTime += frameTime[i]
		}

		write(0, Y_RES-4, "+--------------+")
		write(0, Y_RES-3, fmt.Sprintf("|%9.2f Secs|", time.Now().Sub(TIME_START).Seconds()))
		write(0, Y_RES-2, fmt.Sprintf("|%9.2f FPS |", 1.0/(sumFrameTime/20.0)))
		write(0, Y_RES-1, "+--------------+")

		os.Stdout.Write([]byte(PRINT_CURSOR_POS))
		bytes := []byte{}
		found := false
		for {
			r, err := term.ReadRune()
			if err != nil || r == 'R' {
				break
			}

			if r == '[' {
				found = true
				continue
			}

			if !found {
				continue
			}

			bytes = append(bytes, byte(r))
		}
		rawMouse := string(bytes)

		mousePos := strings.Split(strings.ToUpper(rawMouse), "M")
		lower := strings.Index(rawMouse, "M") < 0

		if len(mousePos) > 0 {
			mousePosArr := strings.Split(mousePos[0], ";")
			mousePosArr[0] = strings.Trim(mousePosArr[0], "<")

			if len(mousePosArr) == 3 {
				mouseX, _ = strconv.Atoi(mousePosArr[1])
				mouseY, _ = strconv.Atoi(mousePosArr[2])

				mouseX -= 2
				mouseY -= 2

				if mousePosArr[0] == "35" {
					lastEvent = fmt.Sprintf("MOUSE MOVED TO %d,%d", mouseX, mouseY)
				}

				if mousePosArr[0] == "0" {
					if lower {
						lastEvent = fmt.Sprintf("MOUSE UP ON %d,%d", mouseX, mouseY)
					} else {
						lastEvent = fmt.Sprintf("MOUSE CLICK ON %d,%d", mouseX, mouseY)
					}
				}
			}
		}

		write(0, Y_RES-5, lastEvent)

		drawCircle(mouseX, mouseY, 10, 5.0)

		clear()
		draw()

		nextFrameTarget := baseTime.Add(SLEEP_DURATION)
		durationSleep := nextFrameTarget.Sub(time.Now())

		if durationSleep > 0 {
			time.Sleep(durationSleep - (time.Millisecond / 2))
		}

		for time.Now().Before(nextFrameTarget) {
		}

		baseTime = time.Now().Add(nextFrameTarget.Sub(time.Now()))
		frameTime[counter] = float64(time.Now().Sub(lastTime).Seconds())
		lastTime = time.Now()

		counter++
		counter %= 20
	}
}
