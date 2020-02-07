package yenc

/*// #cgo LDFLAGS: -static
// #include <unistd.h>
import "C"
*/
import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"io"
	"math"
	"os"
	"runtime"
	"sync"
)

const FILEBUFFER_SIZE = 4194304

var (
	sem = semaphore.NewWeighted(1)
	wg  sync.WaitGroup
)

type yEncReader struct {
	sourceString         string
	sourceLength         int
	currentLine          int
	lineLength           int
	lineIndex            int
	currentIndex         int
	ctx                  context.Context
	legacy               bool
	result               []*bytes.Buffer
	EOF                  bool
	lastEncodedBuffer    int
	PartSize             int
	partBuffer           bytes.Buffer
	finishedBuffers      []*bool
	finishedBuffersMutex sync.Mutex
	currentFragment      int
	outputMutex          sync.Mutex
}

/*
	NewyEnc creates a new yEncReader providing the ReadLine function.
*/
func NewyEnc(sourceString string, lineLength int, legacy bool) *yEncReader {
	return &yEncReader{
		sourceString:         sourceString,
		sourceLength:         len(sourceString),
		lineLength:           lineLength,
		ctx:                  context.TODO(),
		legacy:               legacy,
		lastEncodedBuffer:    -1,
		PartSize:             100000,
		currentFragment:      1,
		outputMutex:          sync.Mutex{},
		finishedBuffersMutex: sync.Mutex{},
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
/
func (encoder *yEncReader) ReadLine() ([]byte, error) {

}*/

func (encoder *yEncReader) encodeFileFragment(targetBuffer *bytes.Buffer, offset int64, sem *semaphore.Weighted) {
	defer sem.Release(1)
	defer wg.Done()

	//fmt.Println(offset)

	file, _ := os.Open(encoder.sourceString)
	defer file.Close()
	_, err := file.Seek(offset*FILEBUFFER_SIZE, 0)
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(file)

	buff := make([]byte, 4096)

	// 4 * 1024 * 1024 = 4 MiB
	for i := 0; i < 1024; i++ {
		bte, err := r.Read(buff)
		if err != nil {
			if err == io.EOF {
				encoder.EOF = true
				break
			}
			panic(err)
		}

		for _, b := range buff[0:bte] {

			encoded, escape := yEncify(b, encoder.legacy)

			if escape {
				targetBuffer.WriteRune('=')
			}
			targetBuffer.WriteByte(encoded)
		}
	}
	encoder.finishedBuffersMutex.Lock()
	*encoder.finishedBuffers[offset] = true
	encoder.finishedBuffersMutex.Unlock()
	go encoder.printFragment()
}

func (encoder *yEncReader) EncodeFile() {
	f, err := os.Open(encoder.sourceString)
	if err != nil {
		panic(err)
	}

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	fragments := fi.Size() / FILEBUFFER_SIZE
	limit := getGoRoutineCount(int(math.Floor(float64(fragments)) + 1))

	ctx := context.TODO()
	sem = semaphore.NewWeighted(int64(limit))
	counter := 0

	for !encoder.EOF {
		if err := sem.Acquire(ctx, 1); err != nil {
			panic(err)
		}
		newBool := false
		encoder.finishedBuffersMutex.Lock()
		encoder.finishedBuffers = append(encoder.finishedBuffers, &newBool)
		encoder.finishedBuffersMutex.Unlock()
		encoder.result = append(encoder.result, bytes.NewBuffer([]byte("")))

		wg.Add(1)
		go encoder.encodeFileFragment(encoder.result[counter], int64(counter), sem)
		counter++
	}
	wg.Wait()
}

func getGoRoutineCount(fragments int) int {
	totalRAM := C.sysconf(C._SC_PHYS_PAGES) * C.sysconf(C._SC_PAGE_SIZE)
	RAMLimit := 128
	if totalRAM < 500*1024*1024 {
		RAMLimit = 1
	} else if totalRAM < 2*1024*1024*1024 && fragments > 32 {
		RAMLimit = 32
	} else if totalRAM < 4*1024*1024*1024 && fragments > 64 {
		RAMLimit = 64
	} else if totalRAM < 16*1024*1024*1024 {
		RAMLimit = 128
	}

	CPULimit := runtime.NumCPU() * 15

	if RAMLimit < CPULimit {
		return RAMLimit
	} else {
		return CPULimit
	}
}

func (encoder *yEncReader) printFragment() {
	encoder.outputMutex.Lock()
	defer encoder.outputMutex.Unlock()

	curBuf := encoder.lastEncodedBuffer + 1
	dataBuffer := make([]byte, 4096)
	fragmentBuffer := make([]byte, encoder.PartSize)

	for *encoder.finishedBuffers[curBuf] {
		EOB := false
		for !EOB {
			_, err := encoder.result[curBuf].Read(dataBuffer)
			if err == io.EOF {
				EOB = true
			} else if err != nil {
				panic(err)
			}
			encoder.partBuffer.Write(dataBuffer)
		}
		encoder.result[curBuf].Reset()
		encoder.lastEncodedBuffer++
		curBuf++
	}

	EOB := false
	for encoder.partBuffer.Len() > encoder.PartSize || encoder.EOF {

		n, err := encoder.partBuffer.Read(fragmentBuffer)
		if err == io.EOF {
			EOB = true
		} else if err != nil {
			panic(err)
		}

		fmt.Printf("=ybegin part=%d line=%d size=%d name=%s\n", encoder.currentFragment, encoder.lineLength, n, encoder.sourceString)
		startIndex := 0

		for true {
			endIndex := startIndex + encoder.lineLength
			if endIndex >= len(fragmentBuffer) {
				endIndex = n - 1
			}
			if fragmentBuffer[endIndex] == byte(61) {
				endIndex++
			}

			fmt.Println(string(fragmentBuffer[startIndex:endIndex]))
			startIndex = endIndex + 1
			if startIndex >= n {
				break
			}
		}

		if encoder.EOF && EOB {
			break
		}
	}
}
