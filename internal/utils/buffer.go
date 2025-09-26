package utils


func GrowBuffer(buf []byte) []byte {
	biggerBuf := make([]byte, len(buf) * 2)
	copy(biggerBuf, buf)
	buf = biggerBuf
	return buf
}

func ShiftBuffer(buf []byte, offset int) []byte {
	shiftBuf := make([]byte, len(buf))
	copy(shiftBuf, buf[offset:])
	buf = shiftBuf
	return buf
}