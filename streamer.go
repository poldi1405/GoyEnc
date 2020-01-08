package yenc

import (
	//"fmt"
	"io"
)

type yEncReader struct {
	sourceString string
	currentIndex int
	lineLength   int
	lineIndex    int
}

func NewyEncReader(sourceString string, lineLength int) *yEncReader {
	return &yEncReader{
		sourceString: sourceString,
		lineLength: lineLength - 1,
	}
}

func yEncify(r byte) (byte, bool) {
	escape := false
	temp := int(r)
	
	// Hex + 42d
	temp += 42

	// % 256d
	temp &= 255

	// if 00h 0Ah 0Dh or 3Dh
	if temp == 0 || temp == 10 || temp == 13 || temp == 61 {
		// + 64d
		temp += 64
		
		// % 256d
		temp &= 255
		escape = true
	}

	return byte(temp), escape
}

func (encoder *yEncReader) Read(p []byte) (int, bool, error) {
	if encoder.currentIndex >= len(encoder.sourceString) {
		return 0, true, io.EOF
	}

	x := len(encoder.sourceString) - encoder.currentIndex
	n, bound := 0, 0
	if x >= len(p) {
		bound = len(p)
	} else if x <= len(p) {
		bound = x
	}

	buf := make([]byte, bound*2)
	for n < bound {
		char, escape := yEncify(encoder.sourceString[encoder.currentIndex])
		if escape {
			buf[n] = '='
			n++
			buf[n] = char
		} else {
			buf[n] = char
		}
		if encoder.lineIndex >= encoder.lineLength {
			return n, true, nil
		}
		n++
		encoder.currentIndex++
	}
	copy(p, buf)
	return n, false, nil
}