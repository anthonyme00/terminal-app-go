package output

import (
	"os"
)

type Output interface {
	Write(output []byte)
}

type OutputANSIStdOut struct {
	f *os.File
}

func (o *OutputANSIStdOut) Write(output []byte) {
	o.f.Write([]byte("\033[0;0H"))
	o.f.Write(output)
}

func NewStdOutput() Output {
	return &OutputANSIStdOut{
		f: os.Stdout,
	}
}
