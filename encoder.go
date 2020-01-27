package yenc

import (
	"context"
	"golang.org/x/sync/semaphore"
	"io"
)

var (
	sem = semaphore.NewWeighted(1)
)

type yEncReader struct {
	sourceString []byte
	sourceLength int
	currentLine  int
	lineLength   int
	lineIndex    int
	currentIndex int
	ctx          context.Context
	legacy       bool
}

/*
	NewyEnc creates a new yEncReader providing the ReadLine function.
*/
func NewyEnc(sourceString []byte, lineLength int, legacy bool) *yEncReader {
	return &yEncReader{
		sourceString: sourceString,
		sourceLength: len(sourceString),
		lineLength:   lineLength,
		ctx:          context.TODO(),
		legacy:       legacy,
	}
}

// TODO: implement file-reader returning yEnc instead of default bytes

func yEncify(r byte, legacy bool) (byte, bool) {
	escape := false
	temp := int(r)

	// bin + 42d
	temp += 42

	// % 256d
	temp %= 256

	// if 00h 0Ah 0Dh or 3Dh and if legacy 09h
	tab := false
	if temp == 9 && legacy {
		tab = true
	}
	if temp == 0 || temp == 10 || temp == 13 || temp == 61 || tab {
		// + 64d
		temp += 64

		// % 256d
		temp %= 256
		escape = true
	}
	return byte(temp), escape
}

/*
	Returns the next Line of yEnc encoded content.
*/
func (encoder *yEncReader) ReadLine() ([]byte, error) {
	if err := sem.Acquire(encoder.ctx, 1); err != nil {
		return nil, err
	}
	defer sem.Release(1)

	var currentMapIndex int
	resultMap := make([]byte, encoder.lineLength+1)

	for currentMapIndex < encoder.lineLength {
		if encoder.currentIndex == encoder.sourceLength {
			return resultMap, io.EOF
		}
		resByte, escape := yEncify(encoder.sourceString[encoder.currentIndex], encoder.legacy)

		if escape {
			resultMap[currentMapIndex] = '='
			currentMapIndex++
		}
		resultMap[currentMapIndex] = resByte
		currentMapIndex++
		encoder.currentIndex++
	}

	return resultMap, nil
}
