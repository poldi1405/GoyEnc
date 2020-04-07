package yenc

import (
	"context"
	"golang.org/x/sync/semaphore"
	"io"
	"os"
)

func (y yEncer) EncodeFile(FilePath string, ResultChannel chan []byte) (int, error) {
	defer close(ResultChannel)

	// Open File
	encFile, err := os.Open(FilePath)
	if err != nil {
		return -1, err
	}
	defer encFile.Close()

	// Get the file's size
	encFileInfo, err := encFile.Stat()
	if err != nil {
		return -1, err
	}
	filesize := encFileInfo.Size()

	// calculate how many chunks there are
	chunks := filesize / y.chunkSize
	if filesize%y.chunkSize != 0 {
		chunks++
	}

	workID := len(y.ctx)
	worker := getLimit(chunks, y.chunkSize)
	workContext, cancel := context.WithCancel(context.Background())
	y.ctx = append(y.ctx, workContext)
	y.ctxCancel = append(y.ctxCancel, cancel)
	encsem := semaphore.NewWeighted(worker)

	enc := func(ctx context.Context, fragment []byte) <-chan []byte {
		output := make(chan []byte)
		go func() {
			defer encsem.Release(1)
			select {
			case <-ctx.Done():
				return
			default:
				result := []byte{}

				for _, b := range fragment {
					eb, esc := y.yEncify(b)
					if esc {
						result = append(result, '=')
					}
					result = append(result, eb)
				}
				output <- result
			}
		}()
		return output
	}

	for i := int64(0); i < chunks; i++ {
		err = encsem.Acquire(workContext, 1)
		if err != nil {
			cancel()
			return -1, err
		}
		chunk := make([]byte, y.chunkSize)

		_, err = encFile.ReadAt(chunk, i*y.chunkSize)
		if err != nil && err != io.EOF {
			cancel()
			return -1, err
		}

		enc(workContext, chunk)
		//TODO: Join Results
		//TODO: make joiner cancellable (is that a word?)
	}

	return workID, nil
}

func (y yEncer) EncodeBytes(fragment []byte) []byte {
	result := []byte{}

	for _, b := range fragment {
		eb, esc := y.yEncify(b)
		if esc {
			result = append(result, '=')
		}
		result = append(result, eb)
	}
	return result
}
