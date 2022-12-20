package mandelbrot

import (
	"math"
	"time"

	"github.com/anthonyme00/terminal-test/window"
)

type complex struct {
	r float64
	i float64
}

const (
	DEG_2_RAD = math.Pi / 180.0
)

func (c *complex) Add(c1 complex) complex {
	return complex{
		r: c.r + c1.r,
		i: c.i + c1.i,
	}
}

func (c *complex) Squared() *complex {
	return &complex{
		r: (c.r * c.r) - (c.i * c.i),
		i: (2 * c.r * c.i),
	}
}

func (c *complex) DistSqrd() float64 {
	return (c.r * c.r) + (c.i * c.i)
}

func Sample(x, y float64, iteration int) float64 {
	c := complex{
		r: x,
		i: y,
	}

	zn := complex{
		r: 0,
		i: 0,
	}

	n := 0

	for n = 0; n < iteration; n++ {
		n_zn := zn.Squared().Add(c)

		if n_zn.DistSqrd() > 4.0 {
			break
		}

		zn = n_zn
	}

	return float64(n) / float64(iteration-1)
}

type Mandelbrot struct {
	ZoomPeriod     float64
	ZoomCenterX    float64
	ZoomCenterY    float64
	ZoomMin        float64
	ZoomMax        float64
	RotationMinDeg float64
	RotationMaxDeg float64

	window    *window.Window
	startTime time.Time
}

func NewMandelbrotApp() window.Application {
	return &Mandelbrot{
		ZoomPeriod:     3.00000,
		ZoomCenterX:    -0.7457,
		ZoomCenterY:    0.11270,
		ZoomMin:        2.00000,
		ZoomMax:        0.00001,
		RotationMinDeg: 0,
		RotationMaxDeg: 0,
	}

}

func (m *Mandelbrot) Init(window *window.Window) {
	m.window = window
	m.startTime = time.Now()
}

func (m *Mandelbrot) Step(u window.UpdatesInfo) {
	xRes, yRes := m.window.GetSize()
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

	absoluteTime := u.Time_AbsoluteTime.Seconds()

	t := math.Pow(math.Sin(math.Pi*(absoluteTime-mandelbrotPeriod/2.0)/mandelbrotPeriod), 4)

	aspect_ratio := float64(xRes) / float64(yRes)
	x_width := lerp(m.ZoomMax, m.ZoomMin, t)
	y_width := x_width / aspect_ratio * terminal_font_y_correction_factor

	getCoord := func(x, y int) (float64, float64) {
		rel_x := float64(x) / float64(xRes-1)
		rel_y := float64(y) / float64(yRes-1)

		x_start, x_end := center_x-x_width, center_x+x_width
		y_start, y_end := center_y-y_width, center_y+y_width

		return lerp(x_start, x_end, rel_x), lerp(y_start, y_end, rel_y)
	}

	for x := 0; x < xRes; x++ {
		for y := 0; y < yRes; y++ {
			s_x, s_y := getCoord(x, y)

			s_x, s_y = rotateAround(center_x, center_y, s_x, s_y, lerp(float64(m.RotationMinDeg), float64(m.RotationMaxDeg), t))

			s := Sample(s_x, s_y, 100)
			i := int(s * (float64(len(graphic)-1) + 0.1))
			m.window.SetScreen(x, y, graphic[i])
		}
	}
}
