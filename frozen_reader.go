package fileutil

import (
	"io"
	"math"
	"time"
)

type FrozenReader struct{}

func (reader FrozenReader) Read(p []byte) (n int, err error) {
	time.Sleep(math.MaxInt32 * time.Second)
	return 0, io.EOF
}
