package output

import (
	"os"
	"os/exec"
)

type OutputANSIStdOut struct {
	f *os.File
}

func (o *OutputANSIStdOut) Write(output []byte) {
	o.f.Write([]byte("\033[0;0H"))
	o.f.Write(output)
}

func (o *OutputANSIStdOut) Open() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (o *OutputANSIStdOut) Close() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func NewStdOutput() Output {
	return &OutputANSIStdOut{
		f: os.Stdout,
	}
}
