package output

type Output interface {
	Open()
	Write(output []byte)
	Close()
}
