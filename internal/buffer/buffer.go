package buffer


type Buf struct {
	B []byte //buf
	R int //index to where the buffer is field up
}

func (buf *Buf) Buffer() []byte {
    return buf.B[:buf.R]
}

func New(capacity int) *Buf {
	if capacity < 8 {
		capacity = 8
	}
	return &Buf{R: 0, B: make([]byte, capacity)}
}

func (buf *Buf) Grow() {
	biggerBuf := make([]byte, len(buf.B) * 2)
	copy(biggerBuf, buf.B)
	buf.B = biggerBuf
}

func (buf *Buf) Free(offset int) {
	shiftBuf := make([]byte, len(buf.B))
	copy(shiftBuf, buf.B[offset:])
	buf.B = shiftBuf
}
