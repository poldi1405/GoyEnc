package yenc

import (
	//"fmt"
	"context"
	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"
	"io"
	"os"
)

type yEncer struct {
	resultSet           []chan []byte
	intermediateResults [][]chan []byte
	finishedResult      []bool
	chunkSize           int64
	experimental        bool // use yEnc 1.3
	lineLength          int

	ctx       []context.Context
	ctxCancel []context.CancelFunc
}

func NewyEnc() yEncer {
	return yEncer{}
}

func NewCustomyEnc(lineLength int, FileChunkSize int64, use13 bool) yEncer {
	return yEncer{
		chunkSize:    FileChunkSize,
		experimental: use13,
		lineLength:   lineLength,
	}
}

func (y yEncer) yEncify(r byte) (byte, bool) {
	escape := false
	temp := int(r)

	// bin + 42d
	temp += 42

	// % 256d
	temp %= 256

	// if 00h 0Ah 0Dh or 3Dh
	if temp == 0 || temp == 10 || temp == 13 || temp == 61 {
		// + 64d
		temp += 64

		// % 256d
		temp %= 256
		escape = true
	}

	return byte(temp), escape
}

func (y yEncer) CancelEncoding(JobId int) error {
	if len(y.ctx) <= JobId {
		return errors.New("yEnc: JobID not assigned")
	} else {
		y.ctxCancel[JobId]()
	}
}
