package general

type Sampler interface {
	Sample(x float64, y float64, params ...interface{})
}
