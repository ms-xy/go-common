package io

type ByteWriter interface {
	Write([]byte) (int, error)
	Get() ([]byte, error)
	MustWrite([]byte) int
	MustGet() []byte
}

type byteWriterImpl struct {
	bytes       []byte
	pos, length int
}

func NewByteWriter() ByteWriter {
	length := 256
	bytes := make([]byte, length)
	writer := &byteWriterImpl{
		bytes:  bytes,
		pos:    0,
		length: length,
	}
	return writer
}

func (this *byteWriterImpl) Write(bytes []byte) (int, error) {
	length := len(bytes)
	if this.pos+length > this.length {
		newBuffer := make([]byte, this.length*2)
		copy(newBuffer, this.bytes)
		this.bytes = newBuffer
	}
	copy(this.bytes[this.pos:], bytes)
	this.pos += length
	return length, nil
}

func (this *byteWriterImpl) Get() ([]byte, error) {
	return this.bytes[:this.pos], nil
}

func (this *byteWriterImpl) MustWrite(bytes []byte) int {
	if n, err := this.Write(bytes); err != nil {
		panic(err)
	} else {
		return n
	}
}

func (this *byteWriterImpl) MustGet() []byte {
	if bytes, err := this.Get(); err != nil {
		panic(err)
	} else {
		return bytes
	}
}
