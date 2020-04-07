package yenc

import (
	"github.com/pbnjay/memory"
)

func getLimit(chunks int64, chunksize int64) int64 {
	// Let's just hope that no one has more than 8191 PeB of RAM
	RAM := int64(memory.TotalMemory())
	if RAM/5 < chunks*chunksize {
		return RAM / 5 / chunksize
	}

	return chunks
}
