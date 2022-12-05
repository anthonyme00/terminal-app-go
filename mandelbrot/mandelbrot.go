package mandelbrot

type complex struct {
	r float64
	i float64
}

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
